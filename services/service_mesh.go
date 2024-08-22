package services

import (
	"context"
	"log/slog"

	"github.com/alekLukanen/go-templ-htmx-example-app/database/queries"
)

type ServiceMesh struct {
	ctx    context.Context
	logger *slog.Logger

	UserAuthenticationService UserAuthenticationService
}

func NewServiceMesh(ctx context.Context, logger *slog.Logger, dbQueries *queries.Queries) ServiceMesh {
	return ServiceMesh{
		ctx:                       ctx,
		logger:                    logger,
		UserAuthenticationService: NewUserAuthenticationService(ctx, logger, dbQueries),
	}
}
