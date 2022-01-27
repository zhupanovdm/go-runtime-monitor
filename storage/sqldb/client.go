package sqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

const dbStorageName = "SQL DB storage"
const defaultTimeout = 5 * time.Second

var _ storage.Storage = (*client)(nil)

type (
	client struct {
		Driver
		db         *sql.DB
		timeout    time.Duration
		dataSource string
		statements map[Query]*sql.Stmt
	}

	Driver interface {
		open(dataSource string) (*sql.DB, error)
		prepare(ctx context.Context, db *sql.DB) error
	}

	Query string
)

const (
	CreateQuery        Query = "INSERT INTO metrics (metric_id, metric_type, value, delta) VALUES ($1,$2,$3,$4)"
	ReadQuery          Query = "SELECT value, delta FROM metrics WHERE metric_id=$1 AND metric_type=$2 LIMIT 1"
	UpdateGaugeQuery   Query = "UPDATE metrics SET value=$3 WHERE metric_id=$1 AND metric_type=$2"
	UpdateCounterQuery Query = "UPDATE metrics SET delta=delta+$3 WHERE metric_id=$1 AND metric_type=$2"
	ReadAllQuery       Query = "SELECT metric_id, metric_type, value, delta FROM metrics"
	DeleteAllQuery     Query = "DELETE FROM metrics"
)

func (c *client) Clear(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	if err := c.queryWithTx(ctx, DeleteAllQuery, exec()); err != nil {
		logger.Err(err).Msg("failed to clear metrics table")
		return err
	}
	logger.Info().Msg("cleared")
	return nil
}

func (c *client) IsPersistent() bool {
	return true
}

func (c *client) Init(ctx context.Context) error {
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName))

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	db, err := c.open(c.dataSource)
	if err != nil {
		logger.Err(err).Msg("failed to open db")
		return err
	}

	// init script may be sensitive to sql dialect and depends on db server
	if err = c.prepare(ctx, db); err != nil {
		logger.Err(err).Msg("failed to prepare db")
		return err
	}

	statements, err := prepareStmts(ctx, db,
		CreateQuery, ReadQuery, UpdateGaugeQuery, UpdateCounterQuery,
		ReadAllQuery, DeleteAllQuery)

	if err != nil {
		logger.Err(err).Msg("failed to prepare statements")
		return err
	}
	c.statements = statements
	c.db = db

	logger.Info().Msg("initialized")
	return nil
}

func (c *client) GetAll(ctx context.Context) (metric.List, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	list := make(metric.List, 0)
	fetchMetrics := fetch(func(_ context.Context, rows *sql.Rows) error {
		m := &Metrics{}
		if err := rows.Scan(&m.ID, &m.typ, &m.value, &m.delta); err != nil {
			return err
		}
		mtr := m.ToCanonical()
		if mtr == nil {
			return errors.New("unable to convert row to canonical metric")
		}
		list = append(list, mtr)
		return nil
	})
	if err := c.queryWithTx(ctx, ReadAllQuery, fetchMetrics); err != nil {
		logger.Err(err).Msg("failed to query all metrics")
		return nil, err
	}

	logger.Trace().Msgf("%d records read", len(list))
	return list, nil
}

func (c *client) UpdateBulk(ctx context.Context, list metric.List) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	return c.withTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		for _, mtr := range list {
			if err := c.process(logging.SetLogger(ctx, logger), tx, mtr); err != nil {
				return err
			}
		}
		logger.Trace().Msgf("%d metrics processed", len(list))
		return nil
	})
}

func (c *client) Get(ctx context.Context, id string, typ metric.Type) (*metric.Metric, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	var mtr *metric.Metric
	fetch := func(ctx context.Context, tx *sql.Tx) (err error) {
		mtr, err = c.read(logging.SetLogger(ctx, logger), tx, id, typ)
		return
	}
	if err := c.withTx(ctx, fetch); err != nil {
		return nil, err
	}
	return mtr, nil
}

func (c *client) Update(ctx context.Context, mtr *metric.Metric) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	return c.withTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		return c.process(logging.SetLogger(ctx, logger), tx, mtr)
	})
}

func (c *client) Ping(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if err := c.db.PingContext(ctx); err != nil {
		logger.Err(err).Msg("failed to ping db server")
		return err
	}

	logger.Trace().Msg("storage is online")
	return nil
}

func (c *client) Close(ctx context.Context) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	if err := c.db.Close(); err != nil {
		logger.Err(err).Msg("failed to close db connection")
	}
	logger.Info().Msg("closed")
}

func (c *client) process(ctx context.Context, tx *sql.Tx, mtr *metric.Metric) error {
	m, err := c.read(ctx, tx, mtr.ID, mtr.Type())
	if err != nil {
		return err
	}
	if m == nil {
		return c.create(ctx, tx, mtr)
	}
	return c.update(ctx, tx, mtr)
}

func (c *client) create(ctx context.Context, tx *sql.Tx, mtr *metric.Metric) error {
	_, logger := logging.GetOrCreateLogger(ctx)
	logger.UpdateContext(logging.LogCtxFrom(mtr))

	if err := mtr.Type().Validate(); err != nil {
		logger.Err(err).Msgf("create failed")
		return err
	}

	stmt, err := c.stmt(ctx, tx, CreateQuery)
	if err != nil {
		logger.Err(err).Msgf("create failed")
		return err
	}
	v, d := toPrimitive(mtr)
	if _, err := stmt.ExecContext(ctx,
		mtr.ID,
		string(mtr.Type()),
		v,
		d); err != nil {
		logger.Err(err).Msgf("create failed")
		return err
	}

	logger.Trace().Msg("created")
	return nil
}

func (c *client) read(ctx context.Context, tx *sql.Tx, id string, typ metric.Type) (*metric.Metric, error) {
	_, logger := logging.GetOrCreateLogger(ctx)
	logger.UpdateContext(logging.LogCtxKeyStr(logging.MetricIDKey, id))
	logger.UpdateContext(logging.LogCtxFrom(typ))

	if err := typ.Validate(); err != nil {
		logger.Err(err).Msgf("query failed")
		return nil, err
	}

	stmt, err := c.stmt(ctx, tx, ReadQuery)
	if err != nil {
		logger.Err(err).Msgf("query failed")
		return nil, err
	}

	var v float64
	var d int64
	if err := stmt.QueryRowContext(ctx,
		id,
		string(typ)).Scan(&v, &d); err != nil {
		if err == sql.ErrNoRows {
			logger.Trace().Msg("not found")
			return nil, nil
		}
		logger.Err(err).Msgf("query failed")
		return nil, err
	}

	var mtr *metric.Metric
	switch typ {
	case metric.GaugeType:
		mtr = metric.NewGaugeMetric(id, metric.Gauge(v))
	case metric.CounterType:
		mtr = metric.NewCounterMetric(id, metric.Counter(d))
	default:
		err := fmt.Errorf("unknown metric %v", typ)
		logger.Err(err).Msgf("query failed")
		return nil, err
	}

	logger.UpdateContext(logging.LogCtxFrom(mtr))
	logger.Trace().Msg("got record")
	return mtr, nil
}

func (c *client) update(ctx context.Context, tx *sql.Tx, mtr *metric.Metric) (err error) {
	_, logger := logging.GetOrCreateLogger(ctx)
	logger.UpdateContext(logging.LogCtxFrom(mtr))

	if err = mtr.Type().Validate(); err != nil {
		logger.Err(err).Msgf("update failed")
		return
	}

	var stmt *sql.Stmt
	v, d := toPrimitive(mtr)
	switch mtr.Type() {
	case metric.GaugeType:
		if stmt, err = c.stmt(ctx, tx, UpdateGaugeQuery); err == nil {
			_, err = stmt.ExecContext(ctx,
				mtr.ID,
				string(mtr.Type()),
				v)
		}
	case metric.CounterType:
		if stmt, err = c.stmt(ctx, tx, UpdateCounterQuery); err == nil {
			_, err = stmt.ExecContext(ctx,
				mtr.ID,
				string(mtr.Type()),
				d)
		}
	default:
		err = fmt.Errorf("unknown metric %v", mtr.Type())
	}
	if err != nil {
		logger.Err(err).Msg("update failed")
		return
	}
	logger.Trace().Msg("updated")
	return
}

func prepareStmts(ctx context.Context, db *sql.DB, queries ...Query) (map[Query]*sql.Stmt, error) {
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName))

	statements := make(map[Query]*sql.Stmt)
	var err error
	for _, query := range queries {
		s := string(query)
		var stmt *sql.Stmt
		if stmt, err = db.PrepareContext(ctx, s); err != nil {
			logger.Err(err).Msgf("failed to prepare statement: %s", s)
			return nil, err
		}
		statements[query] = stmt
	}

	logger.Info().Msgf("%d statements prepared", len(statements))
	return statements, nil
}

func (c *client) stmt(ctx context.Context, tx *sql.Tx, query Query) (*sql.Stmt, error) {
	s, ok := c.statements[query]
	if !ok {
		return nil, fmt.Errorf("query is not prepared statement: %s", query)
	}
	return tx.StmtContext(ctx, s), nil
}

func (c *client) withTx(ctx context.Context, underTx func(ctx context.Context, tx *sql.Tx) error) error {
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Err(err).Msg("failed to open transaction")
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer tx.Rollback()
	if err := underTx(ctx, tx); err != nil {
		logger.Err(err).Msg("failed to exec query")
		return err
	}
	if err := tx.Commit(); err != nil {
		logger.Err(err).Msg("failed to commit transaction")
		return err
	}
	return nil
}

func (c *client) queryWithTx(ctx context.Context, query Query, exec func(ctx context.Context, stmt *sql.Stmt) error) error {
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	return c.withTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		stmt, err := c.stmt(ctx, tx, query)
		if err != nil {
			logger.Err(err).Msg("unable to exec statement")
			return err
		}
		return exec(ctx, stmt)
	})
}

func fetch(fetch func(ctx context.Context, rows *sql.Rows) error, args ...interface{}) func(ctx context.Context, stmt *sql.Stmt) error {
	return func(ctx context.Context, stmt *sql.Stmt) error {
		rows, err := stmt.QueryContext(ctx, args...)
		if err != nil {
			return err
		}
		//goland:noinspection GoUnhandledErrorResult
		defer rows.Close()

		for rows.Next() {
			if err := fetch(ctx, rows); err != nil {
				return err
			}
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return nil
	}
}

func exec(args ...interface{}) func(ctx context.Context, stmt *sql.Stmt) error {
	return func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, args...)
		return err
	}
}

func New(driver Driver) storage.Factory {
	return func(cfg *config.Config) storage.Storage {
		if len(cfg.Database) == 0 {
			return nil
		}
		return &client{
			timeout:    defaultTimeout,
			dataSource: cfg.Database,
			Driver:     driver,
		}
	}
}
