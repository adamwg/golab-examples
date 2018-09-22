package fileserver

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type noopReadCloser struct {
	io.Reader
}

func (nrc *noopReadCloser) Close() error { return nil }

func TestServeFileSuccess(t *testing.T) {
	mos := &MockOS{}
	h := &handler{
		os: mos,
	}

	mos.On("Stat", "/tmp/foo.go").
		Return(&fakeStat{false}, nil).
		Once()
	mos.On("Open", "/tmp/foo.go").
		Return(&noopReadCloser{
			bytes.NewBufferString("hello"),
		}, nil).
		Once()

	w := httptest.NewRecorder()
	u, _ := url.Parse("http://localhost:8080/foo.go")
	req := &http.Request{
		URL: u,
	}

	h.serveFile(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "hello", w.Body.String())
}

func TestServeFileEmptyPath(t *testing.T) {
	mos := &MockOS{}
	h := &handler{
		os: mos,
	}

	w := httptest.NewRecorder()
	u, _ := url.Parse("http://localhost:8080/")
	req := &http.Request{
		URL: u,
	}

	h.serveFile(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestServeFileDirectory(t *testing.T) {
	mos := &MockOS{}
	h := &handler{
		os: mos,
	}

	mos.On("Stat", "/tmp/foo.go").
		Return(&fakeStat{true}, nil).
		Once()

	w := httptest.NewRecorder()
	u, _ := url.Parse("http://localhost:8080/foo.go")
	req := &http.Request{
		URL: u,
	}

	h.serveFile(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestServeFileStatErrors(t *testing.T) {
	mos := &MockOS{}
	h := &handler{
		os: mos,
	}

	mos.On("Stat", "/tmp/foo.go").
		Return(nil, errors.New("oops!")).
		Once()

	w := httptest.NewRecorder()
	u, _ := url.Parse("http://localhost:8080/foo.go")
	req := &http.Request{
		URL: u,
	}

	h.serveFile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestServeFileOpenError(t *testing.T) {
	mos := &MockOS{}
	h := &handler{
		os: mos,
	}

	mos.On("Stat", "/tmp/foo.go").
		Return(&fakeStat{false}, nil).
		Once()
	mos.On("Open", "/tmp/foo.go").
		Return(nil, errors.New("oops!")).
		Once()

	w := httptest.NewRecorder()
	u, _ := url.Parse("http://localhost:8080/foo.go")
	req := &http.Request{
		URL: u,
	}

	h.serveFile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestStatOK(t *testing.T) {
	tcs := []struct {
		name         string
		st           os.FileInfo
		err          error
		expectedCode int
	}{
		{
			name:         "not exist",
			st:           nil,
			err:          os.ErrNotExist,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "permission",
			st:           nil,
			err:          os.ErrPermission,
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "internal",
			st:           nil,
			err:          errors.New("oops"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "directory",
			st:           &fakeStat{true},
			err:          nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "OK",
			st:           &fakeStat{false},
			err:          nil,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			code := statOK(tc.st, tc.err)
			assert.Equal(t, tc.expectedCode, code)
		})
	}
}

type fakeStat struct {
	bool
}

func (s *fakeStat) Name() string {
	panic("not implemented")
}

func (s *fakeStat) Size() int64 {
	panic("not implemented")
}

func (s *fakeStat) Mode() os.FileMode {
	panic("not implemented")
}

func (s *fakeStat) ModTime() time.Time {
	panic("not implemented")
}

func (s *fakeStat) IsDir() bool {
	return s.bool
}

func (s *fakeStat) Sys() interface{} {
	panic("not implemented")
}
