package handler

import (
	"github.com/Golden-Rama-Digital/library-core-go/presentation"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	paymentAuthenticator "github.com/tripdeals/library-service.go"
	errorBackend "github.com/tripdeals/payment.service/src/error"
	"github.com/tripdeals/payment.service/src/middleware"
	"github.com/tripdeals/payment.service/src/utils"
	"github.com/tripdeals/promo.service/config"
	"github.com/tripdeals/promo.service/src/dto"
	"github.com/tripdeals/promo.service/src/service"
)

type MyPromoCodeHandler struct {
	validate         *validator.Validate
	promoCodeService service.PromoCodeServiceV1
	cfg              *config.Config
	authenticator    paymentAuthenticator.PaymentInternalTokenAuthenticatorV1
}

func NewMyPromoCodeHandler(promoCodeService service.PromoCodeServiceV1, validate *validator.Validate, cfg *config.Config, authenticator paymentAuthenticator.PaymentInternalTokenAuthenticatorV1) *PromoCodeHandler {
	return &PromoCodeHandler{validate: validate, promoCodeService: promoCodeService, cfg: cfg, authenticator: authenticator}
}

func (h *MyPromoCodeHandler) Routes(g *echo.Group) {
	g = g.Group("/:platform", middleware.JWTAuthMiddlewareWithPlatforms(h.authenticator, h.cfg.PlatformConfig))
	g.GET("/my-list", h.MyList())
	g.GET("/my-list/:code", h.MyDetail())
	g.POST("/apply", h.Apply())
	g.POST("/redeem", h.Redeem())
}

func (h *MyPromoCodeHandler) MyList() echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			ctx = c.Request().Context()
		)

		userID, ok := utils.UserIDFromContext(ctx)
		if !ok || userID == "" {
			return errorBackend.ErrUserIDNotFound
		}

		return nil
	}
}

func (h *MyPromoCodeHandler) MyDetail() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func (h *MyPromoCodeHandler) Apply() echo.HandlerFunc {
	return func(c echo.Context) error {
		// var (
		// 	ctx = c.Request().Context()
		// )

		// userID, ok := utils.UserIDFromContext(ctx)
		// if !ok || userID == "" {
		// 	return errorBackend.ErrUserIDNotFound
		// }

		// config, _ := utils.GetPlatformConfig(ctx)
		// fmt.Printf("Hello: %s", config)
		// if userID == config.Platform {
		// 	return errorBackend.ErrInvalidToken
		// }

		var req dto.ApplyPromoCodeRequest
		if err := c.Bind(&req); err != nil {
			return presentation.ResponseErrValidation(err)
		}

		result, err := h.promoCodeService.Apply(c.Request().Context(), req)
		if err != nil {
			return err
		}
		return presentation.WriteResponseOk(c, result)
	}
}

func (h *MyPromoCodeHandler) Redeem() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
