// Package postgres - работает с базой
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NikitaBarysh/metrics_and_alertinc/config/server"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"

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

// SetMetrics - добавляем метрики в базу
func (p *Postgres) SetMetrics(metric []entity.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("repository: postgres: SetMetric: BegimTX: %w", err)
	}

	for _, v := range metric {
		var err error

		service.Retry(func() error {
			_, err = tx.ExecContext(ctx,
				insertMetric,
				v.ID,
				v.MType,
				v.Delta,
				v.Value,
			)
			return err
		}, 0)

		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return fmt.Errorf("repository: postgres: SetMetric: Rollback: %w", err)
			}
			return fmt.Errorf("repository: postgres: SetMetric: INSERT INTO: %w", err)
		}
	}

	return tx.Commit()
}

// GetMetric - получаем метрику из базы
func (p *Postgres) GetMetric(key string) (entity.Metric, error) { // TODO
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	metric := entity.Metric{}

	service.Retry(func() error {
		row := p.db.QueryRowContext(ctx, getMetric, key)

		return row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)

	}, 0)

	return metric, nil
}

// GetAllMetric - получаем все метрики из базы
func (p *Postgres) GetAllMetric() ([]entity.Metric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, getAllMetric) //TODo
	if err != nil {
		return nil, fmt.Errorf("repository: postgres: GetAllMetric: QueryContext: %w", err)
	}
	defer rows.Close()

	metricSlice := make([]entity.Metric, 0, 35)

	for rows.Next() {
		m := entity.Metric{}
		err := rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		if err != nil {
			return nil, fmt.Errorf("repository: postgres: GetAllMetric: Scan: %w", err)
		}
		metricSlice = append(metricSlice, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return metricSlice, nil
}

// CheckPing - проверяем соединение с базой
func (p *Postgres) CheckPing(ctx context.Context) error {
	err := p.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("CheckPing: %w", err)
	}
	return nil
}
