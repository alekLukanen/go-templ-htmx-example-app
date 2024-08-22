package componentHandlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/alekLukanen/go-templ-htmx-example-app/handlers"
	"github.com/alekLukanen/go-templ-htmx-example-app/services"
	"github.com/alekLukanen/go-templ-htmx-example-app/ui/components"
)

type AuthRefreshComponentHandler struct {
	ctx    context.Context
	logger *slog.Logger

	*services.ServiceMesh
}

func NewAuthRefreshComponentHandler(ctx context.Context, logger *slog.Logger, serviceMesh *services.ServiceMesh) *AuthRefreshComponentHandler {
	return &AuthRefreshComponentHandler{
		ctx:         ctx,
		logger:      logger,
		ServiceMesh: serviceMesh,
	}
}

func (obj *AuthRefreshComponentHandler) RegisterProtectedRoutes(echoHandler *echo.Group) {
	echoHandler.GET("/auth/refresh", obj.RefreshAuth)
}

func (obj *AuthRefreshComponentHandler) FormToParams(echoCtx echo.Context) components.AuthRefreshParams {
	tokenUser := echoCtx.Get("user").(*jwt.Token)
	claims := tokenUser.Claims.(*services.JWTScopeClaims)

	return components.AuthRefreshParams{JWTExp: claims.ExpiresAt.Unix()}
}

func (obj *AuthRefreshComponentHandler) InitialLoadParams(echoCtx echo.Context) components.AuthRefreshParams {
	return obj.FormToParams(echoCtx)
}

func (obj *AuthRefreshComponentHandler) RefreshAuth(echoCtx echo.Context) error {
	user, err := obj.UserAuthenticationService.UserFromEchoContext(echoCtx)
	if err != nil {
		return err
	}

	token, err := obj.UserAuthenticationService.GenerateJWT(user)
	if err != nil {
		return err
	}

	expireTime := time.Now().UTC().Add(services.TOKEN_DURATION)
	echoCtx.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expireTime,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/r/",
		Secure:   false,
	})

	authComponent := components.AuthRefresh(components.AuthRefreshParams{JWTExp: expireTime.Unix()})
	handlers.Render(echoCtx, &authComponent)

	return nil
}
