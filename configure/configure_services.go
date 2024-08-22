package configure

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/alekLukanen/go-templ-htmx-example-app/database"
	"github.com/alekLukanen/go-templ-htmx-example-app/database/queries"
	"github.com/alekLukanen/go-templ-htmx-example-app/services"
)

type ServiceConfiguration struct {
	ctx    context.Context
	logger *slog.Logger

	ServiceMesh services.ServiceMesh
	DB          *pgxpool.Pool
}

func (obj *ServiceConfiguration) Close(ctx context.Context) error {
	obj.logger.Info("closing core service")

	obj.DB.Close()

	obj.logger.Info("closed core service successfully")

	return nil
}

func NewServiceConfiguration(ctx context.Context, logger *slog.Logger) (*ServiceConfiguration, error) {

	logger.Info("configuring core service")

	innerCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	db, err := database.ConfigurePostgresDBConnection(innerCtx)
	if err != nil {
		return nil, err
	}

	dbQueries := queries.New(db)
	serviceMesh := services.NewServiceMesh(ctx, logger, dbQueries)
	configuration := ServiceConfiguration{
		ctx:         ctx,
		logger:      logger,
		ServiceMesh: serviceMesh,
		DB:          db,
	}
	logger.Info("configured core service successfully")

	return &configuration, nil

}
