package configure

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
)

type HandlerConfigurationFunc func(ctx context.Context, echoHandler *echo.Echo, logger *slog.Logger, configuration *ServiceConfiguration)
