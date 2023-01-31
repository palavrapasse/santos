package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/palavrapasse/santos/internal/logging"
)

const (
	queryServerHostEnvKey = "query_server_host"
	queryServerPortEnvKey = "query_server_port"
)

var (
	queryServerHost               = os.Getenv(queryServerHostEnvKey)
	queryServerPort               = os.Getenv(queryServerPortEnvKey)
	queryServiceURL               = fmt.Sprintf("%s:%s", queryServerHost, queryServerPort)
	queryServiceLeaksEndpoint     = "/leaks"
	queryServicePlatformsEndpoint = "/platforms"
)

func GetLeaks(query string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s?%s", queryServiceURL, queryServiceLeaksEndpoint, query)
	return httpGetQueryService(url)
}

func GetPlatforms(query string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s?%s", queryServiceURL, queryServicePlatformsEndpoint, query)
	return httpGetQueryService(url)
}

func httpGetQueryService(url string) (*http.Response, error) {

	logging.Aspirador.Info(fmt.Sprintf("Calling Query Service: %s", url))

	resp, err := http.Get(url)

	if err != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while calling Query Service: %s", err))
		return nil, err
	}

	logging.Aspirador.Info(fmt.Sprintf("Query Service response status: %d", resp.StatusCode))

	return resp, nil
}
