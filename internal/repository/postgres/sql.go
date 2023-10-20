package postgres

const (
	insertMetric = `INSERT INTO metric (id, "type", delta, "value")
			VALUES($1, $2, $3 ,$4) 
		    ON CONFLICT(id) DO 
		    UPDATE SET delta = metric.delta + excluded.delta ,"value" = excluded.value
`

	getMetric = `
SELECT id, type, delta, value 
FROM metric
WHERE id = $1;
`
	getAllMetric = `
	SELECT *
	FROM metric
`
)
