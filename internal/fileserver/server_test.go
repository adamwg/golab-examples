package fileserver

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPathSuccess(t *testing.T) {
	u, _ := url.Parse("http://localhost:8080/foo.go")
	req := &http.Request{
		URL: u,
	}

	path, err := getPath(req)

	assert.Equal(t, "/tmp/foo.go", path)
	assert.NoError(t, err)
}

func TestGetPathEmpty(t *testing.T) {
	u, _ := url.Parse("http://localhost:8080/")
	req := &http.Request{
		URL: u,
	}

	_, err := getPath(req)

	assert.Error(t, err)
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

func TestRespondSuccess(t *testing.T) {
	r := bytes.NewBufferString("hello")
	w := httptest.NewRecorder()

	respondSuccess(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "hello", w.Body.String())
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
