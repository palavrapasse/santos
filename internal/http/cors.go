package http

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4/middleware"
	"github.com/palavrapasse/santos/internal/logging"
)

const (
	allowOriginsEnvKey = "allow_origins"
)

var (
	allowOrigins = os.Getenv(allowOriginsEnvKey)
)

func CORSConfig() middleware.CORSConfig {
	cc := middleware.DefaultCORSConfig

	if len(allowOrigins) > 0 {
		cc.AllowOrigins = []string{allowOrigins}
	}

	logging.Aspirador.Trace(fmt.Sprintf("yooo: %s | %d", allowOrigins, len(allowOrigins)))
	logging.Aspirador.Trace(fmt.Sprintf("yooo2: %v", cc))

	return cc
}
