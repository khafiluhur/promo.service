package handler

import (
	"strconv"

	"github.com/Golden-Rama-Digital/library-core-go/presentation"
	"github.com/go-playground/validator/v10"
	"github.com/harryosmar/generic-gorm/base"
	"github.com/labstack/echo/v4"
	presentationBackend "github.com/tripdeals/cms.backend.tripdeals.id/src/presentation"
	paymentAuthenticator "github.com/tripdeals/library-service.go"
	errorBackend "github.com/tripdeals/payment.service/src/error"
	"github.com/tripdeals/payment.service/src/middleware"
	"github.com/tripdeals/payment.service/src/utils"
	"github.com/tripdeals/promo.service/config"
	"github.com/tripdeals/promo.service/src/entity"
	"github.com/tripdeals/promo.service/src/service"
)

type PromoCodeHandler struct {
	validate         *validator.Validate
	promoCodeService service.PromoCodeServiceV1
	cfg              *config.Config
	authenticator    paymentAuthenticator.PaymentInternalTokenAuthenticatorV1
}

func NewPromoCodeHandler(promoCodeService service.PromoCodeServiceV1, validate *validator.Validate, cfg *config.Config, authenticator paymentAuthenticator.PaymentInternalTokenAuthenticatorV1) *PromoCodeHandler {
	return &PromoCodeHandler{validate: validate, promoCodeService: promoCodeService, cfg: cfg, authenticator: authenticator}
}

func (p *PromoCodeHandler) Routes(g *echo.Group) {
	g = g.Group("/:platform/promo-code", middleware.JWTAuthMiddlewareWithPlatforms(p.authenticator, p.cfg.PlatformConfig))
	g.GET("", p.List())
	g.GET("/:id", p.Detail())
	g.GET("/by-slugs", p.ByPromoCodes())
	g.POST("", p.Create())
	g.PUT("", p.Update())
	g.PATCH("/activate/:id", p.Activate())
	g.PATCH("/deactivate/:id", p.Deactivate())
}

func (p *PromoCodeHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			defaultPage   = 1
			defaultLimit  = 10
			maxLimit      = 25
			wheres        []base.Where
			orders        []base.OrderBy
			defaultOrders = []base.OrderBy{
				{
					Field:     "updated_at",
					Direction: "desc",
				},
			}
		)

		parseQuery, err := presentationBackend.ParseQuery(
			c.QueryString(),
			entity.PromoCode{}.GetAllowedWhereFields(),
			entity.PromoCode{}.GetAllowedOrderFields(),
		)
		if err == nil {
			wheres = parseQuery.Wheres
			orders = parseQuery.Orders
		}

		pageStr := c.QueryParam("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			page = defaultPage
		}
		if len(orders) == 0 {
			orders = defaultOrders
		}
		limitStr := c.QueryParam("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = defaultLimit
		}
		if limit > maxLimit {
			limit = maxLimit
		}

		ctx := c.Request().Context()
		userID, ok := utils.UserIDFromContext(ctx)
		if !ok || userID == "" {
			return errorBackend.ErrUserIDNotFound
		}

		config, err := utils.GetPlatformConfig(ctx)
		if err != nil {
			return presentation.ResponseErr(errorBackend.ErrConfigPlatform)
		}

		ctx, err = utils.InjectPlatformConfigToContext(ctx, config.Platform, p.cfg.PlatformConfig)
		if err != nil {
			return presentation.ResponseErr(errorBackend.ErrConfigPlatform)
		}
		if userID == config.Platform {
			return errorBackend.ErrInvalidToken
		}

		list, paginator, err := p.promoCodeService.List(ctx, page, limit, orders, wheres)
		if err != nil {
			return err
		}

		return presentation.WritePaging(c, 200, list, &presentation.Paginator{
			Page:    paginator.Page,
			PerPage: paginator.PerPage,
			Total:   int64(paginator.Total),
		})
	}
}

func (p *PromoCodeHandler) Detail() echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()
		userID, ok := utils.UserIDFromContext(ctx)
		if !ok || userID == "" {
			return errorBackend.ErrUserIDNotFound
		}

		config, err := utils.GetPlatformConfig(ctx)
		if err != nil {
			return presentation.ResponseErr(errorBackend.ErrConfigPlatform)
		}

		ctx, err = utils.InjectPlatformConfigToContext(ctx, config.Platform, p.cfg.PlatformConfig)
		if err != nil {
			return presentation.ResponseErr(errorBackend.ErrConfigPlatform)
		}

		if userID == config.Platform {
			return errorBackend.ErrInvalidToken
		}

		result, err := p.promoCodeService.Detail(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return presentation.WriteResponseOk(c, result)
	}
}

func (p *PromoCodeHandler) ByPromoCodes() echo.HandlerFunc {
	return func(c echo.Context) error {
		params := c.QueryParams()
		promoCodes := params["promoCodes[]"]

		ctx := c.Request().Context()
		userID, ok := utils.UserIDFromContext(ctx)
		if !ok || userID == "" {
			return errorBackend.ErrUserIDNotFound
		}

		config, err := utils.GetPlatformConfig(ctx)
		if err != nil {
			return presentation.ResponseErr(errorBackend.ErrConfigPlatform)
		}

		ctx, err = utils.InjectPlatformConfigToContext(ctx, config.Platform, p.cfg.PlatformConfig)
		if err != nil {
			return presentation.ResponseErr(errorBackend.ErrConfigPlatform)
		}

		if userID == config.Platform {
			return errorBackend.ErrInvalidToken
		}

		result, err := p.promoCodeService.ByPromoCodes(c.Request().Context(), promoCodes)
		if err != nil {
			return err
		}

		return presentation.WriteResponseOk(c, result)
	}
}

func (p *PromoCodeHandler) Create() func(c echo.Context) error {
	return func(c echo.Context) error {
		var promoCodeDTO entity.PromoCode
		err := c.Bind(&promoCodeDTO)
		if err != nil {
			return err
		}

		err = p.validate.Struct(promoCodeDTO)
		if err != nil {
			return presentation.ResponseErrValidation(err)
		}

		ctx := c.Request().Context()
		userID, ok := utils.UserIDFromContext(ctx)
		if !ok || userID == "" {
			return errorBackend.ErrUserIDNotFound
		}

		config, err := utils.GetPlatformConfig(ctx)
		if err != nil {
			return presentation.ResponseErr(errorBackend.ErrConfigPlatform)
		}

		ctx, err = utils.InjectPlatformConfigToContext(ctx, config.Platform, p.cfg.PlatformConfig)
		if err != nil {
			return presentation.ResponseErr(errorBackend.ErrConfigPlatform)
		}

		if userID == config.Platform {
			return errorBackend.ErrInvalidToken
		}

		result, err := p.promoCodeService.Create(c.Request().Context(), &promoCodeDTO)
		if err != nil {
			return err
		}
		return presentation.WriteResponseOk(c, result)
	}
}

func (p *PromoCodeHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var promoCodeDTO entity.PromoCode
		err := c.Bind(&promoCodeDTO)
		if err != nil {
			return err
		}

		err = p.validate.Struct(promoCodeDTO)
		if err != nil {
			return presentation.ResponseErrValidation(err)
		}

		ctx := c.Request().Context()
		userID, ok := utils.UserIDFromContext(ctx)
		if !ok || userID == "" {
			return errorBackend.ErrUserIDNotFound
		}

		config, err := utils.GetPlatformConfig(ctx)
		if err != nil {
			return presentation.ResponseErr(errorBackend.ErrConfigPlatform)
		}

		ctx, err = utils.InjectPlatformConfigToContext(ctx, config.Platform, p.cfg.PlatformConfig)
		if err != nil {
			return presentation.ResponseErr(errorBackend.ErrConfigPlatform)
		}

		if userID == config.Platform {
			return errorBackend.ErrInvalidToken
		}

		err = p.promoCodeService.Update(ctx, &promoCodeDTO)
		if err != nil {
			return err
		}
		return presentation.WriteResponseOk(c, nil)
	}
}

func (p *PromoCodeHandler) Activate() echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()
		userID, ok := utils.UserIDFromContext(ctx)
		if !ok || userID == "" {
			return errorBackend.ErrUserIDNotFound
		}

		config, err := utils.GetPlatformConfig(ctx)
		if err != nil {
			return presentation.ResponseErr(errorBackend.ErrConfigPlatform)
		}

		ctx, err = utils.InjectPlatformConfigToContext(ctx, config.Platform, p.cfg.PlatformConfig)
		if err != nil {
			return presentation.ResponseErr(errorBackend.ErrConfigPlatform)
		}

		if userID == config.Platform {
			return errorBackend.ErrInvalidToken
		}

		err = p.promoCodeService.Activate(ctx, id)
		if err != nil {
			return err
		}

		return presentation.WriteResponseOk(c, nil)
	}
}

func (p *PromoCodeHandler) Deactivate() echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()
		userID, ok := utils.UserIDFromContext(ctx)
		if !ok || userID == "" {
			return errorBackend.ErrUserIDNotFound
		}

		config, err := utils.GetPlatformConfig(ctx)
		if err != nil {
			return presentation.ResponseErr(errorBackend.ErrConfigPlatform)
		}

		ctx, err = utils.InjectPlatformConfigToContext(ctx, config.Platform, p.cfg.PlatformConfig)
		if err != nil {
			return presentation.ResponseErr(errorBackend.ErrConfigPlatform)
		}

		if userID == config.Platform {
			return errorBackend.ErrInvalidToken
		}

		err = p.promoCodeService.Deactivate(ctx, id)
		if err != nil {
			return err
		}

		return presentation.WriteResponseOk(c, nil)
	}
}
