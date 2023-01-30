package http

import (
	"fmt"
	"os"
)

const (
	queryServerHostEnvKey = "query_server_host"
	queryServerPortEnvKey = "query_server_port"
)

var (
	queryServerHost = os.Getenv(queryServerHostEnvKey)
	queryServerPort = os.Getenv(queryServerPortEnvKey)
	QueryServiceURL = fmt.Sprintf("%s:%s", queryServerHost, queryServerPort)
	QueryServiceLeakPath = "/leaks"
)
