package handler

import (
	"strconv"

	"github.com/Golden-Rama-Digital/library-core-go/presentation"
	"github.com/go-playground/validator/v10"
	"github.com/harryosmar/generic-gorm/base"
	"github.com/labstack/echo/v4"
	presentationBackend "github.com/tripdeals/cms.backend.tripdeals.id/src/presentation"
	paymentAuthenticator "github.com/tripdeals/library-service.go"
	"github.com/tripdeals/payment.service/src/middleware"
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

func (h *PromoCodeHandler) Routes(g *echo.Group) {
	g = g.Group("/:platform/promo-code", middleware.JWTAuthMiddlewareWithPlatforms(h.authenticator, h.cfg.PlatformConfig))
	g.GET("", h.List())
	g.GET("/:id", h.Detail())
	g.GET("/by-slugs", h.BySlugs())
	g.POST("", h.Create())
	g.PUT("", h.Update())
	g.PATCH("/activate/:id", h.Activate())
	g.PATCH("/deactivate/:id", h.Deactivate())
}

func (h *PromoCodeHandler) List() echo.HandlerFunc {
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
		list, paginator, err := h.promoCodeService.List(ctx, page, limit, orders, wheres)
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

func (h *PromoCodeHandler) Detail() echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}

		result, err := h.promoCodeService.Detail(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return presentation.WriteResponseOk(c, result)
	}
}

func (h *PromoCodeHandler) BySlugs() echo.HandlerFunc {
	return func(c echo.Context) error {
		params := c.QueryParams()
		slugs := params["slugs[]"]
		result, err := h.promoCodeService.BySlugs(c.Request().Context(), slugs)
		if err != nil {
			return err
		}

		return presentation.WriteResponseOk(c, result)
	}
}

func (h *PromoCodeHandler) Create() func(c echo.Context) error {
	return func(c echo.Context) error {
		var promoCodeDTO entity.PromoCode
		err := c.Bind(&promoCodeDTO)
		if err != nil {
			return err
		}

		err = h.validate.Struct(promoCodeDTO)
		if err != nil {
			return presentation.ResponseErrValidation(err)
		}

		result, err := h.promoCodeService.Create(c.Request().Context(), &promoCodeDTO)
		if err != nil {
			return err
		}
		return presentation.WriteResponseOk(c, result)
	}
}

func (h *PromoCodeHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var promoCodeDTO entity.PromoCode
		err := c.Bind(&promoCodeDTO)
		if err != nil {
			return err
		}

		err = h.validate.Struct(promoCodeDTO)
		if err != nil {
			return presentation.ResponseErrValidation(err)
		}

		ctx := c.Request().Context()
		err = h.promoCodeService.Update(ctx, &promoCodeDTO)
		if err != nil {
			return err
		}
		return presentation.WriteResponseOk(c, nil)
	}
}

func (h *PromoCodeHandler) Activate() echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()
		err = h.promoCodeService.Activate(ctx, id)
		if err != nil {
			return err
		}

		return presentation.WriteResponseOk(c, nil)
	}
}

func (h *PromoCodeHandler) Deactivate() echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()
		err = h.promoCodeService.Deactivate(ctx, id)
		if err != nil {
			return err
		}

		return presentation.WriteResponseOk(c, nil)
	}
}
