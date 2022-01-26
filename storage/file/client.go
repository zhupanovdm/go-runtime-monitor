package file

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"sync"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

const fileStorageName = "File storage"

var _ storage.Storage = (*client)(nil)

type client struct {
	sync.RWMutex
	filename string
}

func (c *client) IsPersistent() bool {
	return true
}

func (c *client) Init(context.Context) error {
	return nil
}

func (c *client) Clear(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))
	ctx = logging.SetLogger(ctx, logger)

	if err := os.Remove(c.filename); err != nil && err != fs.ErrNotExist {
		logger.Err(err).Msg("failed to clear destination")
		return err
	}
	return nil
}

func (c *client) GetAll(ctx context.Context) (metric.List, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))
	ctx = logging.SetLogger(ctx, logger)

	logger.Trace().Msg("retrieving metrics from file storage")

	c.RLock()
	defer c.RUnlock()

	r, err := NewJSONFileReader(ctx, c.filename)
	if err != nil {
		logger.Err(err).Msg("file store: failed to open storage for reading")
		return nil, err
	}
	defer r.Close()
	return r.Read()
}

func (c *client) UpdateBulk(ctx context.Context, list metric.List) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))
	ctx = logging.SetLogger(ctx, logger)

	logger.Trace().Msg("persisting metrics to file storage")

	c.Lock()
	defer c.Unlock()

	w, err := NewJSONFileWriter(ctx, c.filename)
	if err != nil {
		return err
	}
	defer w.Close()
	return w.Write(list)
}

func (c *client) Get(ctx context.Context, _ string, _ metric.Type) (*metric.Metric, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))

	logger.Error().Msg("getting metric from storage is unsupported")
	return nil, errors.New("unsupported operation")
}

func (c *client) Update(ctx context.Context, _ string, _ metric.Value) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))

	logger.Error().Msg("update metric in storage is unsupported")
	return errors.New("unsupported operation")
}

func (c *client) Ping(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))

	logger.Trace().Msg("check file availability")

	c.RLock()
	defer c.RUnlock()

	file, err := os.OpenFile(c.filename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		logger.Err(err).Msg("file store: failed to check file availability")
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Err(err).Msg("file store: failed to close file")
		}
	}()
	return nil
}

func (c *client) Close(context.Context) {}

func New(cfg *config.Config) storage.Storage {
	if len(cfg.StoreFile) == 0 {
		return nil
	}
	return &client{filename: cfg.StoreFile}
}
