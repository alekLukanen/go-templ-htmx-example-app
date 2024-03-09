package database

import (
	"errors"
)

var DATABASE_DB_CONNECTION_NOT_CONFIGURED error = errors.New("database connection not configured")
