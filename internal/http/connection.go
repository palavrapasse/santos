package http

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/acme/autocert"
)

const (
	serverHostEnvKey                = "server_host"
	serverPortEnvKey                = "server_port"
	serverDomainEnvKey              = "server_domain"
	serverAutoCertDirFilePathEnvKey = "server_autocrt_cache_dir_fp"
)

var (
	serverHost                      = os.Getenv(serverHostEnvKey)
	serverPort                      = os.Getenv(serverPortEnvKey)
	serverDomain                    = os.Getenv(serverDomainEnvKey)
	serverAutoCertDirectoryFilePath = os.Getenv(serverAutoCertDirFilePathEnvKey)
)

func Start(e *echo.Echo) error {
	addr := serverAddress()

	if len(serverDomain) > 0 {
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(serverDomain)
		e.AutoTLSManager.Cache = autocert.DirCache(serverAutoCertDirectoryFilePath)

		return e.StartAutoTLS(addr)
	}

	return e.Start(addr)
}

func serverAddress() string {
	return fmt.Sprintf("%s:%s", serverHost, serverPort)
}
