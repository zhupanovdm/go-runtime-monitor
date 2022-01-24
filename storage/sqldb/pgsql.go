package sqldb

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var _ Driver = (*PGX)(nil)

type PGX struct{}

func (P PGX) open(dataSource string) (*sql.DB, error) {
	return sql.Open("pgx", dataSource)
}

func (P PGX) init(ctx context.Context, db *sql.DB) (err error) {
	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS metrics (`+
			`metric_id varchar(255) NOT NULL,`+
			`metric_type varchar(255) NOT NULL,`+
			`value double precision,`+
			`delta int8`+
			`);`)
	return
}
