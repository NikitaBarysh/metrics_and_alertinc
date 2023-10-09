package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/config/server"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

type Postgres struct {
	db *sql.DB
}

func InitPostgres(cfg *server.Config) (*Postgres, error) {
	dsn := cfg.DataBaseDSN
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	err = goose.Up(db, ".")
	if err != nil {
		return nil, fmt.Errorf("repository: postgres: InitPostgres: gooose.UP: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("repository: database: CheckConn: %w", err)
	}

	return &Postgres{db: db}, nil
}

func (p *Postgres) GetMetric(key string) map[string]entity.Metric {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := p.db.QueryRowContext(ctx, getMetric, key)

	var metric map[string]entity.Metric

	err := row.Scan(metric[key].ID, metric[key].MType, metric[key].Delta, metric[key].Value)
	if err != nil {
		fmt.Println(fmt.Errorf("repository: postgres: Get: Scan: %w", err))
		return nil
	}

	return metric
}

func (p *Postgres) GetAllMetric() []entity.Metric {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, getAllMetric)
	if err != nil {
		fmt.Println(fmt.Errorf("repository: postgres: GetAllMetric: QueryContext: %w", err))
		return nil
	}

	defer rows.Close()

	metricSlice := make([]entity.Metric, 0, 35)

	for rows.Next() {
		m := entity.Metric{}
		err := rows.Scan(m.ID, m.MType, m.Delta, m.Value)
		if err != nil {
			fmt.Println(fmt.Errorf("repository: postgres: GetAllMetric: Scan: %w", err))
			return nil
		}
		metricSlice = append(metricSlice, m)
	}
	return metricSlice
}

func (p *Postgres) CheckPing(ctx context.Context) error {
	err := p.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("CheckPing: %w", err)
	}
	return nil
}
