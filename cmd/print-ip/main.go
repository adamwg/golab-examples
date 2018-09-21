package main

import (
	"fmt"
	"net/http"

	"github.com/adamwg/golab-examples/internal/ip"
)

func main() {
	ip, err := ip.GetIP(http.DefaultClient)
	fmt.Printf("%s\n", ip)
	fmt.Printf("%v\n", err)
}
