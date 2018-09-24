package time

import (
	"net/http"
	"time"
)

var (
	now       = time.Now
	since     = time.Since
	parseTime = time.Parse
)

func serveTime(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" || path == "/" {
		serveNow(w)
	} else {
		serveSince(w, path[1:])
	}
}

func serveNow(w http.ResponseWriter) {
	t := now().Format(time.RFC3339) + "\n"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(t))
}

func serveSince(w http.ResponseWriter, from string) {
	t, err := parseTime(time.RFC3339, from)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d := since(t).String() + "\n"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(d))
}

func NewHandler() http.Handler {
	return http.HandlerFunc(serveTime)
}
