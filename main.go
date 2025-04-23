package main

import (
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	"github.com/v-escobar/game-api-go/internal/interfaces/common"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/v-escobar/game-api-go/internal/infrastructure/config"

	_ "github.com/v-escobar/game-api-go/docs"
)

type App struct {
	echo *echo.Echo
}

func (a *App) Start(port uint) {
	err := a.echo.Start(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
}

func NewApp(handlers []common.ApplicationHandler) *App {
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fmt.Printf("REQUEST: uri: %v, status: %v\n", v.URI, v.Status)
			return nil
		},
	}))

	for _, handler := range handlers {
		handler.BindRoutes(e)
	}

	return &App{echo: e}
}

// @title			Game API
// @version		0.1
// @description	This is a simple game API.
func main() {
	appConfig := config.InitConfig()
	app := initializeApp(appConfig.Dsn())

	app.Start(appConfig.App.Port)
}
