package main

import (
	"fmt"

	"github.com/adamwg/golab-examples/internal/ip"
)

func main() {
	ip, err := ip.GetIP()
	fmt.Printf("%s\n", ip)
	fmt.Printf("%v\n", err)
}
