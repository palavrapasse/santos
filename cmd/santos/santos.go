package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	as "github.com/palavrapasse/aspirador/pkg"
	"github.com/palavrapasse/santos/internal/http"
	"github.com/palavrapasse/santos/internal/logging"
)

func main() {

	logging.Aspirador = as.WithClients(logging.CreateAspiradorClients(http.ServerAddress()))

	logging.Aspirador.Trace("Starting Santos Service")

	e := echo.New()

	defer e.Close()

	http.RegisterMiddlewares(e)
	http.RegisterHandlers(e)

	e.Logger.Fatal(http.Start(e))
}
