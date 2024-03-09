package configure

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/alekLukanen/go-templ-htmx-example-app/core/settings"
)

const PORT int = 3000

type HttpConfiguration struct {
	ctx    context.Context
	logger *slog.Logger
	Server *http.Server

	ServiceConfiguration *ServiceConfiguration
}

func NewUIHttpConfiguration(ctx context.Context, logger *slog.Logger, serviceConfiguration *ServiceConfiguration) *HttpConfiguration {
	echoHandler := echo.New()

	echoLogger := logger.With(
		slog.String("use", "web"),
		slog.String("framework", "echo"),
	)

	// define middleware
	echoHandler.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				echoLogger.LogAttrs(ctx, slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				echoLogger.LogAttrs(ctx, slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	// basic middlewares
	echoHandler.Use(middleware.Recover())
	echoHandler.Use(middleware.BodyLimit("5M"))
	echoHandler.Use(middleware.Gzip())
	echoHandler.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// rate limiting middleware
	/*
		  config := middleware.RateLimiterConfig{
				Skipper: middleware.DefaultSkipper,
				Store: middleware.NewRateLimiterMemoryStoreWithConfig(
					middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 30, ExpiresIn: 3 * time.Minute},
				),
				IdentifierExtractor: func(ctx echo.Context) (string, error) {
					id := ctx.RealIP()
					return id, nil
				},
				ErrorHandler: func(context echo.Context, err error) error {
					return context.String(http.StatusForbidden, "Unabled to identify client for rate limiting")
				},
				DenyHandler: func(context echo.Context, identifier string, err error) error {
					return context.String(http.StatusTooManyRequests, "Rate limit exceeded")
				},
			}
			echoHandler.Use(middleware.RateLimiterWithConfig(config))
	*/
	// static file server middleware
	staticRootDir := "./core/ui/static"
	if settings.ENVIRONMENT != "local" {
		staticRootDir = "./static"
	}

	staticRoot := echoHandler.Group("/static")
	staticRoot.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   staticRootDir,
		Browse: false,
	}))

	// configure handlers
	ConfigureUIHandlers(ctx, echoHandler, logger, serviceConfiguration)

	// define server
	address := fmt.Sprintf(":%d", PORT)
	server := &http.Server{
		Addr:    address,
		Handler: echoHandler,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, address, l.Addr().String())
			return ctx
		},
	}

	return &HttpConfiguration{
		ctx:                  ctx,
		logger:               logger,
		Server:               server,
		ServiceConfiguration: serviceConfiguration,
	}
}

func (obj *HttpConfiguration) StartServer() {
	err := obj.Server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		obj.logger.Info("server closed")
	} else if err != nil {
		obj.logger.Error("error listening for server", slog.Any("err", err))
	}
}

func (obj *HttpConfiguration) Close(ctx context.Context) error {
	return obj.Server.Shutdown(ctx)
}
