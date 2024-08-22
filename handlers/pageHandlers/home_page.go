package pageHandlers

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"

	"github.com/alekLukanen/go-templ-htmx-example-app/handlers"
	"github.com/alekLukanen/go-templ-htmx-example-app/ui/pages"
)

type HomePageHandler struct {
	ctx    context.Context
	logger *slog.Logger
}

func NewHomePageHandler(ctx context.Context, logger *slog.Logger) *HomePageHandler {
	return &HomePageHandler{
		ctx:    ctx,
		logger: logger,
	}
}

func (obj *HomePageHandler) RegisterPublicRoutes(echoHandler *echo.Echo) {
	echoHandler.GET("/", obj.BasePage)
}

func (obj *HomePageHandler) BasePage(echoCtx echo.Context) error {
	component := pages.Homepage()
	handlers.Render(echoCtx, &component)
	return nil
}
