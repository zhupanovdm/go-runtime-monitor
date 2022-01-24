package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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

func (c *client) Init(ctx context.Context) (err error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(dbStorageName), logging.WithCID(ctx))

	logger.Trace().Msg("init db storage")

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	c.Lock()
	defer c.Unlock()

	var db *sql.DB
	if db, err = c.open(c.dataSource); err != nil {
		logger.Err(err).Msg("failed to open db")
		return
	}

	// init script may be sensitive to sql dialect and depends on db server
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

	if len(list) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	c.Lock()
	defer c.Unlock()

	if _, err := c.db.ExecContext(ctx, `DELETE FROM metrics`); err != nil {
		logger.Err(err).Msg("failed to clear storage")
		return err
	}

	var b strings.Builder
	args := make([]interface{}, 0, 4*len(list))

	// list is not expected to be too large
	for i, mtr := range list {
		data := FromCanonical(mtr)

		var chunk string
		if b.Len() == 0 {
			chunk = "INSERT INTO metrics (metric_id, metric_type, value, delta) VALUES ($%d,$%d,$%d,$%d)"
		} else {
			chunk = ",($%d,$%d,$%d,$%d)"
		}
		b.WriteString(fmt.Sprintf(chunk, 4*i+1, 4*i+2, 4*i+3, 4*i+4))
		args = append(args, data.ID, data.typ, data.value, data.delta)
	}

	q := b.String()
	if _, err := c.db.ExecContext(ctx, q, args...); err != nil {
		logger.Err(err).Msg("failed to fulfill store")
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

func New(driver Driver) storage.Provider {
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
