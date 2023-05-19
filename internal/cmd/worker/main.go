package main

import (
	"github.com/mchmarny/vul/internal/server"
)

var (
	// Could be set at build time.
	version = "v0.0.1-default"
)

func main() {
	server.RunWorker(version)
}
