package file

import (
	"context"
	"sync"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

const fileStorageName = "File storage"

var _ storage.Storage = (*fileStorage)(nil)

type fileStorage struct {
	sync.RWMutex
	filename string
}

func (s *fileStorage) GetAll(ctx context.Context) (metric.List, error) {
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

func (s *fileStorage) UpdateBulk(ctx context.Context, list metric.List) error {
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

func NewStorage(cfg *config.Config) storage.Storage {
	if len(cfg.StoreFile) == 0 {
		return nil
	}

	return &fileStorage{filename: cfg.StoreFile}
}
