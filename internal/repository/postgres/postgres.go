package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/config/server"
	"time"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

type Postgres struct {
	db        *sql.DB
	dbStorage DBStorage
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

//func (p *Postgres) SetMetricToDB(metric entity.Metric) {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
//	defer cancel()
//
//}

func (p *Postgres) GetMetricFromDB(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := p.db.QueryRowContext(ctx, getMetric, key)

	metric := p.dbStorage.MetricMap

	err := row.Scan(metric[key].ID, metric[key].MType, metric[key].Delta, metric[key].Value)
	if err != nil {
		fmt.Println(fmt.Errorf("repository: postgres: Get: Scan: %w", err))
	}

	p.dbStorage.SetMetric(metric)
}

func (p *Postgres) GetAllMetricFromDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, getAllMetric)
	if err != nil {
		fmt.Println(fmt.Errorf("repository: postgres: GetAllMetric: QueryContext: %w", err))
	}

	defer rows.Close()

	metricSlice := make([]entity.Metric, 0, 35)
	metric := p.dbStorage.MetricMap

	for rows.Next() {
		m := entity.Metric{}
		err := rows.Scan(m.ID, m.MType, m.Delta, m.Value)
		if err != nil {
			fmt.Println(fmt.Errorf("repository: postgres: GetAllMetric: Scan: %w", err))
		}
	}

	for _, v := range metricSlice {
		metric[v.ID] = entity.Metric{ID: v.ID, MType: v.MType, Delta: v.Delta, Value: v.Value}
	}
	p.dbStorage.SetMetric(metric)
}

func (p *Postgres) CheckPing(ctx context.Context) error {
	err := p.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("CheckPing: %w", err)
	}
	return nil
}
