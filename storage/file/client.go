package file

import (
	"context"
	"errors"
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

func (s *client) IsPersistent() bool {
	return true
}

func (s *client) Init(context.Context) error {
	return nil
}

func (s *client) GetAll(ctx context.Context) (metric.List, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))
	ctx = logging.SetLogger(ctx, logger)

	logger.Trace().Msg("retrieving metrics from file storage")

	s.RLock()
	defer s.RUnlock()

	r, err := NewJSONFileReader(ctx, s.filename)
	if err != nil {
		logger.Err(err).Msg("file store: failed to open storage for reading")
		return nil, err
	}
	defer r.Close()
	return r.Read()
}

func (s *client) UpdateBulk(ctx context.Context, list metric.List) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))
	ctx = logging.SetLogger(ctx, logger)

	logger.Trace().Msg("persisting metrics to file storage")

	s.Lock()
	defer s.Unlock()

	w, err := NewJSONFileWriter(ctx, s.filename)
	if err != nil {
		return err
	}
	defer w.Close()
	return w.Write(list)
}

func (s *client) Get(ctx context.Context, _ string, _ metric.Type) (*metric.Metric, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))

	logger.Error().Msg("getting metric from storage is unsupported")
	return nil, errors.New("unsupported operation")
}

func (s *client) Update(ctx context.Context, _ string, _ metric.Value) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))

	logger.Error().Msg("update metric in storage is unsupported")
	return errors.New("unsupported operation")
}

func (s *client) Ping(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(fileStorageName), logging.WithCID(ctx))

	logger.Trace().Msg("check file availability")

	s.RLock()
	defer s.RUnlock()

	file, err := os.OpenFile(s.filename, os.O_RDWR|os.O_CREATE, 0777)
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

func (s *client) Close(context.Context) {}

func New(cfg *config.Config) storage.Storage {
	if len(cfg.StoreFile) == 0 {
		return nil
	}
	return &client{filename: cfg.StoreFile}
}
