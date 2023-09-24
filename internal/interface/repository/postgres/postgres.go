package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/config/server"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	cfg *server.Config
}

func NewPostgres(config *server.Config) *Postgres {
	return &Postgres{cfg: config}
}

func (p *Postgres) InitPostgres() (*sql.DB, error) {
	//dsn := p.cfg.DataBaseDSN
	db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable")
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("repository: database: CheckConn: %w", err)
	}
	return db, nil
}
