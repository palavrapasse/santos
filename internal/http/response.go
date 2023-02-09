package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func TooManyRequests(ectx echo.Context) error {
	return ectx.NoContent(http.StatusTooManyRequests)
}

func Ok(ectx echo.Context, i interface{}) error {
	return ectx.JSON(http.StatusOK, i)
}

func InternalServerError(ectx echo.Context) error {
	return ectx.NoContent(http.StatusInternalServerError)
}
