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

type StrikeThroughtPriceHandler struct {
	validate                   *validator.Validate
	strikeThroughtPriceService service.StrikeThroughtPriceServiceV1
	cfg                        *config.Config
	authenticator              paymentAuthenticator.PaymentInternalTokenAuthenticatorV1
}

func NewStrikeThroughtPriceHandler(strikeThroughtPriceService service.StrikeThroughtPriceServiceV1, validate *validator.Validate, cfg *config.Config, authenticator paymentAuthenticator.PaymentInternalTokenAuthenticatorV1) *StrikeThroughtPriceHandler {
	return &StrikeThroughtPriceHandler{validate: validate, strikeThroughtPriceService: strikeThroughtPriceService, cfg: cfg, authenticator: authenticator}
}

func (s *StrikeThroughtPriceHandler) Routes(g *echo.Group) {
	g = g.Group("/:platform/strike-throught-price", middleware.JWTAuthMiddlewareWithPlatforms(s.authenticator, s.cfg.PlatformConfig))
	g.GET("", s.List())
	g.GET("/:id", s.Detail())
	g.GET("/by-slugs", s.ByIds())
	g.POST("", s.Create())
	g.PUT("", s.Update())
	g.PATCH("/activate/:id", s.Activate())
	g.PATCH("/deactivate/:id", s.Deactivate())
}

func (s *StrikeThroughtPriceHandler) List() echo.HandlerFunc {
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
		list, paginator, err := s.strikeThroughtPriceService.List(ctx, page, limit, orders, wheres)
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

func (s *StrikeThroughtPriceHandler) Detail() echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}

		result, err := s.strikeThroughtPriceService.Detail(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return presentation.WriteResponseOk(c, result)
	}
}

func (s *StrikeThroughtPriceHandler) ByIds() echo.HandlerFunc {
	return func(c echo.Context) error {
		params := c.QueryParams()
		ids := params["ids[]"]
		result, err := s.strikeThroughtPriceService.BySlugs(c.Request().Context(), ids)
		if err != nil {
			return err
		}

		return presentation.WriteResponseOk(c, result)
	}
}

func (s *StrikeThroughtPriceHandler) Create() func(c echo.Context) error {
	return func(c echo.Context) error {
		var strikeThroughtPriceDTO entity.StrikeThroughtPrice
		err := c.Bind(&strikeThroughtPriceDTO)
		if err != nil {
			return err
		}

		err = s.validate.Struct(strikeThroughtPriceDTO)
		if err != nil {
			return presentation.ResponseErrValidation(err)
		}

		result, err := s.strikeThroughtPriceService.Create(c.Request().Context(), &strikeThroughtPriceDTO)
		if err != nil {
			return err
		}
		return presentation.WriteResponseOk(c, result)
	}
}

func (s *StrikeThroughtPriceHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var strikeThroughtPriceDTO entity.StrikeThroughtPrice
		err := c.Bind(&strikeThroughtPriceDTO)
		if err != nil {
			return err
		}

		err = s.validate.Struct(strikeThroughtPriceDTO)
		if err != nil {
			return presentation.ResponseErrValidation(err)
		}

		ctx := c.Request().Context()
		err = s.strikeThroughtPriceService.Update(ctx, &strikeThroughtPriceDTO)
		if err != nil {
			return err
		}
		return presentation.WriteResponseOk(c, nil)
	}
}

func (s *StrikeThroughtPriceHandler) Activate() echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()
		err = s.strikeThroughtPriceService.Activate(ctx, id)
		if err != nil {
			return err
		}

		return presentation.WriteResponseOk(c, nil)
	}
}

func (s *StrikeThroughtPriceHandler) Deactivate() echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()
		err = s.strikeThroughtPriceService.Deactivate(ctx, id)
		if err != nil {
			return err
		}

		return presentation.WriteResponseOk(c, nil)
	}
}
