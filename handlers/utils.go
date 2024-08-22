package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(echoCtx echo.Context, component *templ.Component) {
	(*component).Render(echoCtx.Request().Context(), echoCtx.Response().Writer)
}
