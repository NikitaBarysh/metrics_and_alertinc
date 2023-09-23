package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	conn *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{conn: db}
}

func (r *Repository) CheckConn(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("repository: database: CheckConn: %w", err)
	}
	return nil
}
