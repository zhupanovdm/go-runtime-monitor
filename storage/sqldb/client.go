package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

const dbStorageName = "SQL DB storage"
const defaultTimeout = 5 * time.Second

var _ storage.Storage = (*client)(nil)

type client struct {
	sync.RWMutex
	db         *sql.DB
	timeout    time.Duration
	dataSource string
	Driver
}

func (c *client) Clear(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	c.RLock()
	defer c.RUnlock()
	if _, err := c.db.ExecContext(ctx, "DELETE FROM metrics"); err != nil {
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
	c.Lock()
	defer c.Unlock()
	if err := c.prepare(ctx, db); err != nil {
		logger.Err(err).Msg("failed to prepare db")
		return err
	}
	c.db = db

	logger.Info().Msg("initialized")
	return nil
}

func (c *client) GetAll(ctx context.Context) (metric.List, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	c.RLock()
	defer c.RUnlock()
	rows, err := c.db.QueryContext(ctx, `SELECT metric_id, metric_type, value, delta FROM metrics`)
	if err != nil {
		logger.Err(err).Msg("failed to query all metrics")
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Err(err).Msg("failed to close query result")
		}
	}()

	list := make(metric.List, 0)
	for rows.Next() {
		m := &Metrics{}
		if err := rows.Scan(&m.ID, &m.typ, &m.value, &m.delta); err != nil {
			logger.Err(err).Msg("failed to read query result row")
			return nil, err
		}
		mtr := m.ToCanonical()
		if mtr == nil {
			logger.Err(err).Msg("unable to convert row to canonical metric")
			return nil, err
		}
		list = append(list, mtr)
	}
	if err := rows.Err(); err != nil {
		logger.Err(err).Msg("malformed query result reading detected")
		return nil, err
	}

	logger.Trace().Msgf("%d records read", len(list))
	return list, nil
}

func (c *client) UpdateBulk(ctx context.Context, list metric.List) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	// naive implementation
	c.Lock()
	defer c.Unlock()
	for _, mtr := range list {
		if err := c.process(logging.SetLogger(ctx, logger), mtr); err != nil {
			return err
		}
	}
	logger.Trace().Msgf("%d metrics processed", len(list))
	return nil
}

func (c *client) Get(ctx context.Context, id string, typ metric.Type) (*metric.Metric, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	c.RLock()
	defer c.RUnlock()
	return c.read(logging.SetLogger(ctx, logger), id, typ)
}

func (c *client) Update(ctx context.Context, mtr *metric.Metric) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	c.Lock()
	defer c.Unlock()
	if err := c.process(logging.SetLogger(ctx, logger), mtr); err != nil {
		return err
	}
	return nil
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

func (c *client) process(ctx context.Context, mtr *metric.Metric) error {
	m, err := c.read(ctx, mtr.ID, mtr.Type())
	if err != nil {
		return err
	}
	if m == nil {
		return c.create(ctx, mtr)
	}
	return c.update(ctx, mtr)
}

func (c *client) create(ctx context.Context, mtr *metric.Metric) error {
	_, logger := logging.GetOrCreateLogger(ctx)
	logger.UpdateContext(logging.LogCtxFrom(mtr))

	if err := mtr.Type().Validate(); err != nil {
		logger.Err(err).Msgf("create failed")
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	v, d := toPrimitive(mtr)
	if _, err := c.db.ExecContext(ctx, "INSERT INTO metrics (metric_id, metric_type, value, delta) VALUES ($1,$2,$3,$4)", mtr.ID, string(mtr.Type()), v, d); err != nil {
		logger.Err(err).Msgf("create failed")
		return err
	}

	logger.Trace().Msg("created")
	return nil
}

func (c *client) read(ctx context.Context, id string, typ metric.Type) (*metric.Metric, error) {
	_, logger := logging.GetOrCreateLogger(ctx)
	logger.UpdateContext(logging.LogCtxKeyStr(logging.MetricIDKey, id))
	logger.UpdateContext(logging.LogCtxFrom(typ))

	if err := typ.Validate(); err != nil {
		logger.Err(err).Msgf("query failed")
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	var v float64
	var d int64
	if err := c.db.QueryRowContext(ctx, `SELECT value, delta FROM metrics WHERE metric_id=$1 AND metric_type=$2 LIMIT 1`, id, string(typ)).Scan(&v, &d); err != nil {
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

func (c *client) update(ctx context.Context, mtr *metric.Metric) (err error) {
	_, logger := logging.GetOrCreateLogger(ctx)
	logger.UpdateContext(logging.LogCtxFrom(mtr))

	if err = mtr.Type().Validate(); err != nil {
		logger.Err(err).Msgf("update failed")
		return
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	v, d := toPrimitive(mtr)
	switch mtr.Type() {
	case metric.GaugeType:
		_, err = c.db.ExecContext(ctx,
			"UPDATE metrics SET value = $1 WHERE metric_id = $2 AND metric_type=$3", v, mtr.ID, string(mtr.Type()))
	case metric.CounterType:
		_, err = c.db.ExecContext(ctx,
			"UPDATE metrics SET delta = delta + $1 WHERE metric_id = $2 AND metric_type=$3", d, mtr.ID, string(mtr.Type()))
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

type Driver interface {
	open(dataSource string) (*sql.DB, error)
	prepare(ctx context.Context, db *sql.DB) error
}
