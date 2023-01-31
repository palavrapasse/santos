package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/palavrapasse/santos/internal/logging"
)

const (
	queryServerHostEnvKey = "query_server_host"
	queryServerPortEnvKey = "query_server_port"
)

var (
	queryServerHost      = os.Getenv(queryServerHostEnvKey)
	queryServerPort      = os.Getenv(queryServerPortEnvKey)
	queryServiceURL      = fmt.Sprintf("%s:%s", queryServerHost, queryServerPort)
	queryServiceLeakPath = "/leaks"
)

func GetLeaks(query string) (interface{}, error) {
	url := fmt.Sprintf("%s%s?%s", queryServiceURL, queryServiceLeakPath, query)

	logging.Aspirador.Info(fmt.Sprintf("Calling Query Service: %s", url))

	resp, err := http.Get(url)

	if err != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while calling Query Service: %s", err))
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while reading body of Query Service response: %s", err))
		return nil, err
	}

	var response interface{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while unmarshal body: %s", err))
		return nil, err
	}

	logging.Aspirador.Trace("Received leaks from Query Service")

	return response, nil
}
