package main

import (
	"net/http"

	"github.com/adamwg/golab-examples/internal/time"
)

func main() {
	handler := time.NewHandler()
	http.ListenAndServe(":8080", handler)
}
