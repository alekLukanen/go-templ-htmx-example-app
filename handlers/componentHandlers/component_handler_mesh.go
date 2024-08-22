package componentHandlers

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/alekLukanen/go-templ-htmx-example-app/handlers"
	"github.com/alekLukanen/go-templ-htmx-example-app/services"
	"github.com/labstack/echo/v4"
)

type ComponentHandlerMesh struct {
	ctx    context.Context
	logger *slog.Logger

	AuthRefreshComponentHandler *AuthRefreshComponentHandler
}

func NewComponentHandlerMesh(ctx context.Context, logger *slog.Logger, serviceMesh *services.ServiceMesh) ComponentHandlerMesh {
	return ComponentHandlerMesh{
		ctx:                         ctx,
		logger:                      logger,
		AuthRefreshComponentHandler: NewAuthRefreshComponentHandler(ctx, logger, serviceMesh),
	}
}

func (obj *ComponentHandlerMesh) RegisterProtectedRoutes(echoHandler *echo.Group) {

	val := reflect.ValueOf(*obj)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Check if the field type implements the specified interface
		if field.Type().Implements(reflect.TypeOf((*handlers.ProtectedUIHandler)(nil)).Elem()) {
			// Convert the field value to the interface type and get its pointer
			field.Interface().(handlers.ProtectedUIHandler).RegisterProtectedRoutes(echoHandler)
		}
	}

}
