package handler

import (
	"github.com/Golden-Rama-Digital/library-core-go/dto"
	"github.com/Golden-Rama-Digital/library-core-go/presentation"
	"github.com/labstack/echo/v4"
)

var HealthCheckHandler = func(appName string, version string) func(c echo.Context) error {
	return func(c echo.Context) error {
		return presentation.WriteResponseOk(c, dto.HealthCheck{
			Message: "OK",
			AppName: appName,
			Version: version,
		})
	}
}
