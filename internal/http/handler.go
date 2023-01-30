package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	url := QueryServiceURL + QueryServiceLeakPath + ectx.QueryString()
	logging.Aspirador.Info(fmt.Sprintf("Calling Query Service: %s", url))

	resp, err := http.Get(url)

	if err != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while calling Query Service: %s", err))
		return InternalServerError(ectx)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while reading body of Query Service response: %s", err))
		return InternalServerError(ectx)
	}

	logging.Aspirador.Trace("Success in querying leaks")

	var data interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while unmarshal body: %s", err))
		return InternalServerError(ectx)
	}

	return Ok(ectx, data)
}

func useNotFoundHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	}
}
