package main

import (
	"fmt"
	"strings"

	"github.com/Golden-Rama-Digital/library-core-go/locales"
	customMiddleware "github.com/Golden-Rama-Digital/library-core-go/middleware"
	"github.com/Golden-Rama-Digital/library-core-go/presentation"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	serviceMiddleware "github.com/tripdeals/payment.service/src/middleware"
	"github.com/tripdeals/payment.service/src/utils"
	"github.com/tripdeals/promo.service/config"
	"github.com/tripdeals/promo.service/src/dto"
	"github.com/tripdeals/promo.service/src/factory"
	"github.com/tripdeals/promo.service/src/handler"
)

var (
	validate         *validator.Validate
	err              error
	cfg              *config.Config
	platformConfigs  map[string]dto.Config
	prometheusConfig echoprometheus.MiddlewareConfig
)

func init() {
	cfg = config.Get()
	validate, err = locales.InitValidate(locales.GetTrans(), cfg.AppLocale)
	if err != nil {
		panic(err)
	}

	monitor()
}

func monitor() {
	containerId := uuid.New().String()
	prometheusConfig = echoprometheus.MiddlewareConfig{
		Skipper: func(c echo.Context) bool {
			path := c.Path()
			hasSuffix := strings.HasSuffix(path, "metrics")
			if hasSuffix {
				return true
			}
			hasSuffix = strings.HasSuffix(path, "favicon.ico")
			return hasSuffix
		},
		LabelFuncs: map[string]echoprometheus.LabelValueFunc{
			"container_id": func(c echo.Context, err error) string { // additional custom label
				return containerId
			},
			"env": func(c echo.Context, err error) string { // additional custom label
				return cfg.AppEnv
			},
			"version": func(c echo.Context, err error) string { // additional custom label
				return cfg.AppVersion
			},
		},
		Subsystem: strings.ReplaceAll(cfg.AppName, "-", "_"),
	}

}

func main() {
	e := echo.New()

	e.Use(echoprometheus.NewMiddlewareWithConfig(prometheusConfig))
	e.GET("/metrics", echoprometheus.NewHandler())

	e.Pre(middleware.RemoveTrailingSlash(), middleware.BodyLimit("2M"))

	platformConfigsBase64 := cfg.PlatformConfig
	configs, err := utils.MustDecodeBase64ToStruct[[]dto.Config](platformConfigsBase64)
	if err != nil {
		panic(fmt.Errorf("failed to decode platform configs: %w", err))
	}
	platformConfigs = make(map[string]dto.Config)
	for _, config := range configs {
		platformConfigs[config.Platform] = config
	}

	e.Debug = cfg.AppDebug
	e.HTTPErrorHandler = presentation.CustomHTTPErrorHandler

	db, err := factory.ResolveDatabase(cfg)
	if err != nil {
		panic(err)
	}

	redisClient := factory.ResolveRedisClient(cfg)
	redisCache := factory.ResolveCache(redisClient)
	defer redisCache.Close()
	paymentAuthenticator := factory.ResolvePaymentAuthenticator(cfg)

	// Initialize all services
	promoCodeService := factory.ResolvePromoCodeService(db)
	strikeThroughtPriceService := factory.ResolveStrikeThroughtPriceService(db)

	e.Use(middleware.Recover(), customMiddleware.CorsMiddleware, customMiddleware.RequestIdMiddlewareV2(cfg.AppEnv, cfg.AppVersion))
	defaultMiddlewares := append(
		[]echo.MiddlewareFunc{},
		customMiddleware.LoggerWithContextMiddlewareV2...,
	)

	e.Use()
	e.Use(middleware.CORS())
	e.GET("/", handler.HealthCheckHandler(cfg.AppName, cfg.AppVersion))
	e.GET(cfg.AppApiPath, handler.HealthCheckHandler(cfg.AppName, cfg.AppVersion))

	paymentGroup := e.Group(cfg.AppApiPath)
	paymentGroup.Use(defaultMiddlewares...)
	paymentGroup.Use(serviceMiddleware.JWTAuthMiddlewareWithPlatforms(paymentAuthenticator, cfg.PlatformConfig))
	handler.NewPromoCodeHandler(promoCodeService, validate, cfg, paymentAuthenticator).Routes(paymentGroup)
	handler.NewStrikeThroughtPriceHandler(strikeThroughtPriceService, validate, cfg, paymentAuthenticator).Routes(paymentGroup)

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", cfg.AppPort)))
}
