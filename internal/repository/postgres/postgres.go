package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/config/server"
	"github.com/jackc/pgx/v5"
	"time"

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

func (p *Postgres) SetMetrics(metric []entity.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	//fmt.Println("1111")

	p.db.C

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("repository: postgres: SetMetric: BegimTX: %w", err)
	}
	//fmt.Println("2222")

	for _, v := range metric {
		//fmt.Println("3333")

		//fmt.Println(v.ID, v.MType, v.Delta, v.Value)
		_, err := tx.ExecContext(ctx, `INSERT INTO metric (id, "type", delta, "value")
			VALUES($1, $2, $3 ,$4) 
		    ON CONFLICT(id) DO 
		    UPDATE SET delta = metric.delta + excluded.delta ,"value" = excluded.value`,
			v.ID,
			v.MType,
			v.Delta,
			v.Value,
		)
		//v.ID, v.MType, v.Delta, v.Value)
		//fmt.Println(v.ID, v.MType, v.Delta, v.Value)
		//fmt.Println("4444")
		if err != nil {
			//fmt.Println("err", err)
			err := tx.Rollback()
			if err != nil {
				return fmt.Errorf("repository: postgres: SetMetric: Rollback: %w", err)
			}
			return fmt.Errorf("repository: postgres: SetMetric: INSERT INTO: %w", err)
		}
		//fmt.Println("5555")
	}
	//fmt.Println("6666")
	return tx.Commit()
}

func (p *Postgres) UpdateGaugeMetric(key string, value float64) {
	metric := entity.Metric{ID: key, MType: "gauge", Value: value, Delta: 0}
	err := p.SetMetrics([]entity.Metric{metric})
	if err != nil {
		fmt.Println(fmt.Errorf("repository: postgres: UpdateGauge: SetMetric: %w", err))
	}
}

func (p *Postgres) UpdateCounterMetric(key string, value int64) {
	//fmt.Println("11")
	metric := entity.Metric{ID: key, MType: "counter", Delta: value, Value: 0}
	//fmt.Println("22")
	err := p.SetMetrics([]entity.Metric{metric})
	if err != nil {
		fmt.Println(fmt.Errorf("repository: postgres: UpdateCounter: SetMetric: %w", err))
	}
	//fmt.Println("33")
}

func (p *Postgres) GetMetric(key string) (entity.Metric, error) { // TODO
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := p.db.QueryRowContext(ctx, getMetric, key)

	metric := entity.Metric{}

	err := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
		return metric, err // TODO
	}
	if err != nil {
		return metric, fmt.Errorf("repository: postgres: Get: Scan: %w", err)
	}

	return metric, nil
}

func (p *Postgres) GetAllMetric() ([]entity.Metric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, getAllMetric)
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

	return metricSlice, nil
}

func (p *Postgres) CheckPing(ctx context.Context) error {
	err := p.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("CheckPing: %w", err)
	}
	return nil
}
