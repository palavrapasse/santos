package http

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/palavrapasse/santos/internal/logging"
)

func RegisterMiddlewares(e *echo.Echo) {
	e.Use(loggingMiddleware())
}

func loggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {

			req := ectx.Request()

			logging.Aspirador.Info(fmt.Sprintf("Host: %s | Method: %s | Path: %s", req.Host, req.Method, req.URL.RequestURI()))

			return next(ectx)
		}
	}
}
