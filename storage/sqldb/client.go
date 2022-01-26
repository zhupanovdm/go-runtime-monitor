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

const dbStorageName = "DB storage"
const defaultTimeout = 5 * time.Second

var _ storage.Storage = (*client)(nil)

type client struct {
	sync.RWMutex
	db         *sql.DB
	timeout    time.Duration
	dataSource string
	Driver
}

func (c *client) IsPersistent() bool {
	return true
}

func (c *client) Init(ctx context.Context) (err error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	logger.Trace().Msg("init db storage")

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	var db *sql.DB
	if db, err = c.open(c.dataSource); err != nil {
		logger.Err(err).Msg("failed to open db")
		return
	}

	// init script may be sensitive to sql dialect and depends on db server
	c.Lock()
	defer c.Unlock()
	if err = c.init(ctx, db); err != nil {
		logger.Err(err).Msg("failed to init db")
		return
	}
	c.db = db
	return
}

func (c *client) GetAll(ctx context.Context) (metric.List, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))
	logger.Trace().Msg("retrieving metrics from db storage")

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
		if err := m.Scan(rows); err != nil {
			logger.Err(err).Msg("failed to read row from query result")
			return nil, err
		}
		mtr := m.ToCanonical()
		if mtr == nil {
			logger.Err(err).Msg("failed to read row from query result")
			return nil, err
		}
		list = append(list, mtr)
	}
	if err := rows.Err(); err != nil {
		logger.Err(err).Msg("malformed query result")
		return nil, err
	}
	return list, nil
}

func (c *client) UpdateBulk(ctx context.Context, list metric.List) error {
	// naive implementation

	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))
	logger.Trace().Msg("persisting metrics to db storage")

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	c.Lock()
	defer c.Unlock()
	for _, mtr := range list {
		if err := c.save(ctx, mtr.ID, mtr.Value); err != nil {
			return err
		}
	}
	logger.Trace().Msgf("%d metrics processed", len(list))
	return nil
}

func (c *client) Get(ctx context.Context, id string, typ metric.Type) (*metric.Metric, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))
	logger.Trace().Msg("getting metric from db storage")

	c.RLock()
	defer c.RUnlock()
	v, err := c.read(ctx, id, typ)
	if err != nil {
		logger.Err(err).Msg("failed to retrieve metric from db")
		return nil, err
	}
	if v != nil {
		return &metric.Metric{ID: id, Value: v}, nil
	}
	return nil, nil
}

func (c *client) Update(ctx context.Context, id string, value metric.Value) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))
	logger.Trace().Msg("updating gauge in db storage")

	c.Lock()
	defer c.Unlock()
	if err := c.save(ctx, id, value); err != nil {
		logger.Err(err).Msg("failed to save metric to db")
		return err
	}
	return nil
}

func (c *client) Ping(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))
	logger.Trace().Msg("ping db server")

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if err := c.db.PingContext(ctx); err != nil {
		logger.Err(err).Msg("failed to ping db server")
		return err
	}
	return nil
}

func (c *client) Close(ctx context.Context) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	if err := c.db.Close(); err != nil {
		logger.Err(err).Msg("failed to close db connection")
	}
}

func (c *client) save(ctx context.Context, id string, value metric.Value) error {
	v, err := c.read(ctx, id, value.Type())
	if err != nil {
		return err
	}
	if v == nil {
		if err := c.create(ctx, id, value); err != nil {
			return err
		}
	}
	return c.update(ctx, id, value)
}

func (c *client) create(ctx context.Context, id string, value metric.Value) error {
	if err := value.Type().Validate(); err != nil {
		return err
	}

	v, d := toPrimitive(value)
	if _, err := c.db.ExecContext(ctx, "INSERT INTO metrics (metric_id, metric_type, value, delta) VALUES ($1,$2,$3,$4)", id, string(value.Type()), v, d); err != nil {
		return err
	}
	return nil
}

func (c *client) read(ctx context.Context, id string, typ metric.Type) (metric.Value, error) {
	var v float64
	var d int64

	if err := c.db.QueryRowContext(ctx, `SELECT value, delta FROM metrics WHERE metric_id=$1 AND metric_type=$2 LIMIT 1`, id, string(typ)).Scan(&v, &d); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	switch typ {
	case metric.GaugeType:
		return metric.Value((*metric.Gauge)(&v)), nil
	case metric.CounterType:
		return metric.Value((*metric.Counter)(&d)), nil
	}
	return nil, fmt.Errorf("unknown metric type %v", typ)
}

func (c *client) update(ctx context.Context, id string, value metric.Value) error {
	if err := value.Type().Validate(); err != nil {
		return err
	}
	v, d := toPrimitive(value)
	switch value.Type() {
	case metric.GaugeType:
		if _, err := c.db.ExecContext(ctx, "UPDATE metrics SET value = $1 WHERE metric_id = $2 AND metric_type=$3", v, id, string(value.Type())); err != nil {
			return err
		}
	case metric.CounterType:
		if _, err := c.db.ExecContext(ctx, "UPDATE metrics SET delta = delta + $1 WHERE metric_id = $2 AND metric_type=$3", d, id, string(value.Type())); err != nil {
			return err
		}
	}
	return nil
}

func toPrimitive(value metric.Value) (v float64, d int64) {
	if m, ok := value.(*metric.Metric); ok {
		value = m.Value
	}
	switch value.Type() {
	case metric.GaugeType:
		v = float64(*value.(*metric.Gauge))
	case metric.CounterType:
		d = int64(*value.(*metric.Counter))
	}
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
	init(ctx context.Context, db *sql.DB) error
}
