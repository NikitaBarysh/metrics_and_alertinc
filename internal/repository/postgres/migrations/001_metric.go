package migrations

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upTable, downTable)
}

func upTable(tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS metric (
    "id" VARCHAR(255) UNIQUE NOT NULL,
    "type" VARCHAR(50) NOT NULL,
    "delta" BIGINT,
    "value" DOUBLE PRECISION);`
	_, err := tx.Exec(query)
	if err != nil {
		return fmt.Errorf("migrations: upTable: %w", err)
	}
	return nil
}

func downTable(tx *sql.Tx) error {
	query := `DROP TABLE metric`
	_, err := tx.Exec(query)
	if err != nil {
		return fmt.Errorf("migrations: downTable: %w", err)
	}
	return nil
}
