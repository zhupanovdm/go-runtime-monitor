package handlers

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var baseURL = "http://localhost:8080"

func TestHandlerDo(t *testing.T) {
	t.Run("Basic test", func(t *testing.T) {
		resp := httptest.NewRecorder()
		var h Handler = func(w http.ResponseWriter, r *http.Request) {
			writeBody(w, "foo")
		}

		h.Do(resp, httptest.NewRequest("POST", baseURL, nil))

		result := resp.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.Equal(t, "foo", resp.Body.String())
	})

	t.Run("Nil handler", func(t *testing.T) {
		var h Handler
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("POST", baseURL, nil)

		result := resp.Result()
		defer result.Body.Close()

		assert.NotPanics(t, func() {
			h.Do(resp, req)
		})
	})
}

func TestPOST(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantNext   bool
		wantStatus int
	}{
		{
			name:       "Basic test",
			method:     "POST",
			wantStatus: http.StatusOK,
			wantNext:   true,
		},
		{
			name:       "GET test",
			method:     "GET",
			wantStatus: http.StatusMethodNotAllowed,
			wantNext:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			nextCalled := false

			POST(resp, httptest.NewRequest(tt.method, baseURL, nil), func(http.ResponseWriter, *http.Request) {
				nextCalled = true
			})

			result := resp.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.wantStatus, result.StatusCode)
			if tt.wantNext {
				assert.True(t, nextCalled, "handler call is expected")
			} else {
				assert.False(t, nextCalled, "handler call is not expected")
			}
		})
	}
}

func TestSequence(t *testing.T) {
	tests := []struct {
		name        string
		middlewares []Middleware
		want        string
		wantStatus  int
	}{
		{
			name: "Basic test",
			middlewares: []Middleware{
				func(w http.ResponseWriter, r *http.Request, h Handler) {
					if writeBody(w, "foo") {
						h.Do(w, r)
					}
				},
			},
			wantStatus: http.StatusOK,
			want:       "foo",
		},
		{
			name: "Multiple middlewares",
			middlewares: []Middleware{
				func(w http.ResponseWriter, r *http.Request, h Handler) {
					if writeBody(w, "1") {
						h.Do(w, r)
					}
				},
				func(w http.ResponseWriter, r *http.Request, h Handler) {
					if writeBody(w, "2") {
						h.Do(w, r)
					}
				},
				func(w http.ResponseWriter, r *http.Request, h Handler) {
					if writeBody(w, "3") {
						h.Do(w, r)
					}
				},
			},
			wantStatus: http.StatusOK,
			want:       "123",
		},
		{
			name: "Cancel further middlewares",
			middlewares: []Middleware{
				func(w http.ResponseWriter, r *http.Request, h Handler) {
					writeBody(w, "1")
				},
				func(w http.ResponseWriter, r *http.Request, h Handler) {
					writeBody(w, "2")
				},
				func(w http.ResponseWriter, r *http.Request, h Handler) {
					writeBody(w, "3")
				},
			},
			wantStatus: http.StatusOK,
			want:       "1",
		},
		{
			name:        "Zero middlewares",
			middlewares: []Middleware{},
			wantStatus:  http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()

			Handle(tt.middlewares...).ServeHTTP(resp, httptest.NewRequest("POST", baseURL, nil))

			result := resp.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.wantStatus, result.StatusCode)
			assert.Equal(t, tt.want, resp.Body.String())
		})
	}
}

func writeBody(w http.ResponseWriter, body string) bool {
	if _, err := w.Write(bytes.NewBufferString(body).Bytes()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	return true
}
