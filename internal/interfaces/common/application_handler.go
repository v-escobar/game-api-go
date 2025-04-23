package common

import "github.com/labstack/echo/v4"

type ApplicationHandler interface {
	BindRoutes(e *echo.Echo)
}
