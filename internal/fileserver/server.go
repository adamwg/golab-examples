package fileserver

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func serveFile(w http.ResponseWriter, r *http.Request) {
	path, err := getPath(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	st, err := os.Stat(path)
	if code := statOK(st, err); code != http.StatusOK {
		w.WriteHeader(code)
		return
	}

	f, err := os.Open(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	respondSuccess(w, f)
}

func getPath(r *http.Request) (string, error) {
	path := r.URL.Path
	if path == "" || path == "/" {
		return "", errors.New("no path in request")
	}
	return filepath.Join("/tmp", path), nil
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

func respondSuccess(w http.ResponseWriter, f io.Reader) {
	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}

func NewHandler() http.Handler {
	return http.HandlerFunc(serveFile)
}
