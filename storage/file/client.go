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

func (c *client) Init(ctx context.Context) error {
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName))
	logger.Info().Msg("initialized")
	return nil
}

func (c *client) Clear(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))

	if err := os.Remove(c.filename); err != nil && err != fs.ErrNotExist {
		logger.Err(err).Msg("unable to remove destination file")
		return err
	}
	logger.Info().Msg("cleared")
	return nil
}

func (c *client) GetAll(ctx context.Context) (metric.List, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))

	c.RLock()
	defer c.RUnlock()
	r, err := NewJSONFileReader(logging.SetLogger(ctx, logger), c.filename)
	if err != nil {
		logger.Err(err).Msg("failed to create reader")
		return nil, err
	}
	defer r.Close()

	list, err := r.Read()
	if err != nil {
		logger.Err(err).Msg("failed to read entire file")
		return nil, err
	}
	return list, err
}

func (c *client) UpdateBulk(ctx context.Context, list metric.List) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))

	c.Lock()
	defer c.Unlock()
	w, err := NewJSONFileWriter(logging.SetLogger(ctx, logger), c.filename)
	if err != nil {
		logger.Err(err).Msg("failed to create writer")
		return err
	}
	defer w.Close()

	if err := w.Write(list); err != nil {
		logger.Err(err).Msg("metrics update failed")
		return err
	}
	logger.Trace().Msgf("%d records updated", len(list))
	return nil
}

func (c *client) Get(ctx context.Context, _ string, _ metric.Type) (*metric.Metric, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))

	logger.Error().Msg("get metric operation is unsupported")
	return nil, errors.New("unsupported operation")
}

func (c *client) Update(ctx context.Context, _ *metric.Metric) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))

	logger.Error().Msg("update metric operation is unsupported")
	return errors.New("unsupported operation")
}

func (c *client) Ping(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))

	c.RLock()
	defer c.RUnlock()
	file, err := os.OpenFile(c.filename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		logger.Err(err).Msg("failed to check destination file availability")
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Err(err).Msg("failed to close destination file")
		}
	}()

	logger.Trace().Msg("storage is online")
	return nil
}

func (c *client) Close(ctx context.Context) {
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName))
	logger.Info().Msg("closed")
}

func New(cfg *config.Config) storage.Storage {
	if len(cfg.StoreFile) == 0 {
		return nil
	}
	return &client{filename: cfg.StoreFile}
}
