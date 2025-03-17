package drawing

import (
	"backend/pkg/connectors/pgconnector"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pgConnector *pgconnector.Connector) *Repository {
	return &Repository{
		pool: pgConnector.GetPool(),
	}
}
