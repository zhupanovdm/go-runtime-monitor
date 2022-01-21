package monitor

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerCompressDecompress(t *testing.T) {
	root := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {})
	server := httptest.NewServer(entryHandler(root, compress, decompress))
	defer server.Close()

	tests := []struct {
		name       string
		header     http.Header
		wantStatus int
		wantHeader http.Header
	}{
		{
			name:       "Basic test",
			header:     http.Header{"Accept-Encoding": {"gzip"}},
			wantStatus: http.StatusOK,
			wantHeader: http.Header{"Content-Encoding": {"gzip"}},
		},
		{
			name:       "No compression",
			wantStatus: http.StatusOK,
			wantHeader: http.Header{"Content-Encoding": {""}},
		},
		{
			name:       "Unsupported compression algorithm",
			header:     http.Header{"Accept-Encoding": {"br"}},
			wantStatus: http.StatusOK,
			wantHeader: http.Header{"Content-Encoding": {""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer([]byte("test")))
			require.NoError(t, err)
			req.Header = tt.header

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			for k := range tt.wantHeader {
				assert.Equal(t, tt.wantHeader.Get(k), resp.Header.Get(k))
			}
		})
	}
}
