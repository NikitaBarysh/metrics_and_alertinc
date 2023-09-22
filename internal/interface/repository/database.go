package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	conn *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{conn: db}
}

func (r *Repository) CheckConn(db *sqlx.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("repository: database: CheckConn: %w", err)
	}
	return nil
}
