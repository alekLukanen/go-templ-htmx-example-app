package pageHandlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alekLukanen/go-templ-htmx-example-app/handlers"
	"github.com/alekLukanen/go-templ-htmx-example-app/handlers/componentHandlers"
	"github.com/alekLukanen/go-templ-htmx-example-app/services"
	"github.com/alekLukanen/go-templ-htmx-example-app/ui/pages"
)

type AppHandler struct {
	ctx    context.Context
	logger *slog.Logger

	*services.ServiceMesh
	*componentHandlers.ComponentHandlerMesh
}

func NewAppHandler(ctx context.Context, logger *slog.Logger, serviceMesh *services.ServiceMesh, componentHandlerMesh *componentHandlers.ComponentHandlerMesh) *AppHandler {
	return &AppHandler{
		ctx:                  ctx,
		logger:               logger,
		ServiceMesh:          serviceMesh,
		ComponentHandlerMesh: componentHandlerMesh,
	}
}

func (obj *AppHandler) RegisterProtectedRoutes(echo *echo.Group) {
	echo.GET("/", obj.BasePage)
}

func (obj *AppHandler) BasePage(echoCtx echo.Context) error {
	user, err := obj.UserAuthenticationService.UserFromEchoContext(echoCtx)
	if err != nil {
		return echoCtx.String(http.StatusUnauthorized, "Unauthorized")
	}
	params := pages.AppParams{
		UserEmail:         user.Email,
		AuthRefreshParams: obj.AuthRefreshComponentHandler.InitialLoadParams(echoCtx),
	}
	page := pages.App(params)

	handlers.Render(echoCtx, &page)
	return nil
}
