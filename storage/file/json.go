package file

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"os"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
)

type jsonWriter struct {
	ctx     context.Context
	dest    io.WriteCloser
	encoder *json.Encoder
}

func (w *jsonWriter) Write(list metric.List) error {
	_, logger := logging.GetOrCreateLogger(w.ctx)

	for _, mtr := range list {
		if err := w.encoder.Encode(mtr); err != nil {
			logger.Err(err).Msgf("json write: failed to encode metric %v", mtr)
			return err
		}
	}

	logger.Trace().Msgf("json write: %d records written", len(list))
	return nil
}

func (w *jsonWriter) Close() {
	if w == nil {
		return
	}
	_, logger := logging.GetOrCreateLogger(w.ctx)

	if err := w.dest.Close(); err != nil {
		logger.Err(err).Msg("json write: failed to close destination")
		return
	}
}

func NewJsonWriter(ctx context.Context, dest io.WriteCloser) *jsonWriter {
	return &jsonWriter{
		ctx:     ctx,
		dest:    dest,
		encoder: json.NewEncoder(dest),
	}
}

func NewJsonFileWriter(ctx context.Context, fileName string) (*jsonWriter, error) {
	_, logger := logging.GetOrCreateLogger(ctx)

	if err := os.Remove(fileName); err != nil && err != fs.ErrNotExist {
		logger.Err(err).Msg("json write: failed to clear destination")
		return nil, err
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		logger.Err(err).Msg("json write: failed to open destination")
		return nil, err
	}

	logger.Trace().Msgf("json write: opened for write: %s", file.Name())
	return NewJsonWriter(ctx, file), nil
}

type jsonReader struct {
	ctx     context.Context
	src     io.ReadCloser
	scanner *bufio.Scanner
}

func (r *jsonReader) Read() (metric.List, error) {
	_, logger := logging.GetOrCreateLogger(r.ctx)

	list := make(metric.List, 0)
	for r.scanner.Scan() {
		mtr := &metric.Metric{}
		list = append(list, mtr)
		data := r.scanner.Bytes()
		if err := json.Unmarshal(data, mtr); err != nil {
			logger.Err(err).Msgf("json read: failed to decode '%s'", string(data))
			return nil, err
		}
	}
	if err := r.scanner.Err(); err != nil {
		logger.Err(err).Msg("json read: failed to read source")
		return nil, err
	}

	logger.Trace().Msgf("json read: %d records read", len(list))
	return list, nil
}

func (r *jsonReader) Close() {
	if r == nil {
		return
	}
	_, logger := logging.GetOrCreateLogger(r.ctx)

	if err := r.src.Close(); err != nil {
		logger.Err(err).Msg("json read: failed to close source")
	}
}

func NewJsonReader(ctx context.Context, src io.ReadCloser) *jsonReader {
	return &jsonReader{
		ctx:     ctx,
		src:     src,
		scanner: bufio.NewScanner(src),
	}
}

func NewJsonFileReader(ctx context.Context, fileName string) (*jsonReader, error) {
	_, logger := logging.GetOrCreateLogger(ctx)
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		logger.Err(err).Msg("json read: failed to open source")
		return nil, err
	}

	logger.Trace().Msgf("json read: opened for read: %s", file.Name())
	return NewJsonReader(ctx, file), nil
}
