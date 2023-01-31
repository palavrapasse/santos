package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/palavrapasse/santos/internal/logging"
)

func RegisterHandlers(e *echo.Echo) {

	e.GET(leaksRoute, QueryLeaks)
	e.GET(platformsRoute, QueryPlatforms)

	echo.NotFoundHandler = useNotFoundHandler()
}

func QueryLeaks(ectx echo.Context) error {
	logging.Aspirador.Trace("Querying leaks")

	response, err := GetLeaks(ectx.QueryString())

	if err != nil {
		return InternalServerError(ectx)
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while reading body of Query Service response: %s", err))
		return InternalServerError(ectx)
	}

	logging.Aspirador.Trace("Success in querying leaks")

	return ectx.JSONBlob(response.StatusCode, body)
}

func QueryPlatforms(ectx echo.Context) error {
	logging.Aspirador.Trace("Querying Platforms")

	response, err := GetPlatforms(ectx.QueryString())

	if err != nil {
		return InternalServerError(ectx)
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while reading body of Query Service response: %s", err))
		return InternalServerError(ectx)
	}

	logging.Aspirador.Trace("Success in querying platforms")

	return ectx.JSONBlob(response.StatusCode, body)
}

func useNotFoundHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	}
}
