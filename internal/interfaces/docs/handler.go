package docs

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type DocumentationHandler struct {
}

func NewDocumentationHandler() *DocumentationHandler {
	return &DocumentationHandler{}
}

func (d *DocumentationHandler) BindRoutes(e *echo.Echo) {
	e.GET("/docs/*", echoSwagger.WrapHandler)
}
