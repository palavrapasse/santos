package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/palavrapasse/santos/internal/logging"
)

func RegisterHandlers(e *echo.Echo) {

	e.GET(leaksRoute, QueryLeaks)

	echo.NotFoundHandler = useNotFoundHandler()
}

func QueryLeaks(ectx echo.Context) error {
	logging.Aspirador.Trace("Querying leaks")

	response, err := GetLeaks(ectx.QueryString())

	if err != nil {
		return InternalServerError(ectx)
	}

	logging.Aspirador.Trace("Success in querying leaks")

	return Ok(ectx, response)
}

func useNotFoundHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	}
}
