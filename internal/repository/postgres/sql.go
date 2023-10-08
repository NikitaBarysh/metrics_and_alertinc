package postgres

const (
	createTable = `
CREATE TABLE IF NOT EXISTS praktikum (
    "id" VARCHAR(255) UNIQUE NOT NULL,
    "type" VARCHAR(50) NOT NULL,
    "delta" BIGINT,
    "value" DOUBLE PRECISION,
)
`
	getMetric = `
SELECT * 
FROM metrics 
WHERE id = $1 AND type = $2;
`
	getAllMetric = `
	SELECT *
	FROM metrics
`
)
