package handlers

import (
	"github.com/labstack/echo/v4"
)

type PublicUIHandler interface {
	RegisterPublicRoutes(*echo.Echo)
}

type ProtectedUIHandler interface {
	RegisterProtectedRoutes(*echo.Group)
}
