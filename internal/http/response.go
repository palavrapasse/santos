package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func TooManyRequests(ectx echo.Context) error {
	return ectx.NoContent(http.StatusTooManyRequests)
}

func InternalServerError(ectx echo.Context) error {
	return ectx.NoContent(http.StatusInternalServerError)
}
