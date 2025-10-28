package handler

import (
	"net/http"

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
	validate           *validator.Validate
	myPromoCodeService service.MyPromoCodeServiceV1
	cfg                *config.Config
	authenticator      paymentAuthenticator.PaymentInternalTokenAuthenticatorV1
}

func NewMyPromoCodeHandler(myPromoCodeService service.MyPromoCodeServiceV1, validate *validator.Validate, cfg *config.Config, authenticator paymentAuthenticator.PaymentInternalTokenAuthenticatorV1) *MyPromoCodeHandler {
	return &MyPromoCodeHandler{validate: validate, myPromoCodeService: myPromoCodeService, cfg: cfg, authenticator: authenticator}
}

func (m *MyPromoCodeHandler) Routes(g *echo.Group) {
	g = g.Group("/:platform", middleware.JWTAuthMiddlewareWithPlatforms(m.authenticator, m.cfg.PlatformConfig))
	g.GET("/my-list", m.MyList())
	g.GET("/my-list/:code", m.MyDetail())
	g.POST("/apply", m.Apply())
	g.POST("/redeem", m.Redeem())
}

func (m *MyPromoCodeHandler) MyList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		userID, ok := utils.UserIDFromContext(ctx)
		if !ok || userID == "" {
			return errorBackend.ErrUserIDNotFound
		}

		config, _ := utils.GetPlatformConfig(ctx)
		if userID == config.Platform {
			return errorBackend.ErrInvalidToken
		}

		list, err := m.myPromoCodeService.MyList(ctx, userID)
		if err != nil {
			return err
		}

		return presentation.WriteResponseOk(c, list)
	}
}

func (m *MyPromoCodeHandler) MyDetail() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		userID, ok := utils.UserIDFromContext(ctx)
		if !ok || userID == "" {
			return errorBackend.ErrUserIDNotFound
		}

		config, _ := utils.GetPlatformConfig(ctx)
		if userID == config.Platform {
			return errorBackend.ErrInvalidToken
		}

		code := c.Param("code")
		if code == "" {
			return presentation.WriteResponseCreated(c, http.StatusBadRequest, "Kode promo tidak boleh kosong")
		}

		data, err := m.myPromoCodeService.MyDetail(ctx, code, userID)
		if err != nil {
			return err
		}

		return presentation.WriteResponseOk(c, data)
	}
}

func (m *MyPromoCodeHandler) Apply() echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			ctx = c.Request().Context()
		)

		userID, ok := utils.UserIDFromContext(ctx)
		if !ok || userID == "" {
			return errorBackend.ErrUserIDNotFound
		}

		config, _ := utils.GetPlatformConfig(ctx)
		if userID == config.Platform {
			return errorBackend.ErrInvalidToken
		}

		var req dto.ApplyPromoCodeRequest
		if err := c.Bind(&req); err != nil {
			return presentation.ResponseErrValidation(err)
		}

		result, err := m.myPromoCodeService.Apply(c.Request().Context(), req)
		if err != nil {
			return err
		}
		return presentation.WriteResponseOk(c, result)
	}
}

func (m *MyPromoCodeHandler) Redeem() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		userID, ok := utils.UserIDFromContext(ctx)
		if !ok || userID == "" {
			return errorBackend.ErrUserIDNotFound
		}

		config, _ := utils.GetPlatformConfig(ctx)
		if userID == config.Platform {
			return errorBackend.ErrInvalidToken
		}

		var req dto.RedeemPromoRequest
		if err := c.Bind(&req); err != nil {
			return presentation.ResponseErrValidation(err)
		}

		if err := m.validate.Struct(req); err != nil {
			return presentation.ResponseErrValidation(err)
		}

		req.UserID = userID
		resp, err := m.myPromoCodeService.Redeem(ctx, req)
		if err != nil {
			return err
		}

		return presentation.WriteResponseOk(c, resp)
	}
}
