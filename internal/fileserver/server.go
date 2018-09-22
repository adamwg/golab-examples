package fileserver

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type handler struct {
	os OS
}

func (h *handler) serveFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" || path == "/" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	path = filepath.Join("/tmp", path)

	st, err := h.os.Stat(path)
	if code := statOK(st, err); code != http.StatusOK {
		w.WriteHeader(code)
		return
	}

	f, err := h.os.Open(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}

func statOK(st os.FileInfo, err error) int {
	switch {
	case os.IsNotExist(err):
		return http.StatusNotFound
	case os.IsPermission(err):
		return http.StatusForbidden
	case err != nil:
		return http.StatusInternalServerError
	case st.IsDir():
		return http.StatusBadRequest
	}

	return http.StatusOK
}

func NewHandler() http.Handler {
	h := &handler{
		os: &realOS{},
	}
	return http.HandlerFunc(h.serveFile)
}
