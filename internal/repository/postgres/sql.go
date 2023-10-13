package postgres

const (
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
