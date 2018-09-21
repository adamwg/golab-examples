package main

import (
	"net/http"

	"github.com/adamwg/golab-examples/internal/fileserver"
)

func main() {
	handler := fileserver.NewHandler()
	http.ListenAndServe(":8080", handler)
}
