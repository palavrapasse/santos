package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func InternalServerError(ectx echo.Context) error {
	return ectx.NoContent(http.StatusInternalServerError)
}
