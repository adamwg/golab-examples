package fileserver

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func serveFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	path = filepath.Join("/tmp", path)
	st, err := os.Stat(path)
	switch {
	case os.IsNotExist(err):
		w.WriteHeader(http.StatusNotFound)
		return
	case os.IsPermission(err):
		w.WriteHeader(http.StatusForbidden)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		return
	case st.IsDir():
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	f, err := os.Open(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}

func NewHandler() http.Handler {
	return http.HandlerFunc(serveFile)
}
