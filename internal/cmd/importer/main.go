package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mchmarny/vul/internal/server"
)

var (
	// Set at build time.
	version = "v0.0.1-default"

	imageURI string
	filePath string
	connStr  string
	logLevel string
)

func init() {
	flag.StringVar(&imageURI, "image", "", "The URI of the image that was used to generate the report.")
	flag.StringVar(&filePath, "file", "", "The path of the vulnerability report.")
	flag.StringVar(&connStr, "conn", "", "Target DB connection string.")
	flag.StringVar(&logLevel, "log", "info", "Log level (debug, info, warn, error, fatal, panic)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func usage() {
	flag.Usage()
	os.Exit(1)
}

func main() {
	flag.Parse()

	if imageURI == "" || filePath == "" || connStr == "" {
		usage()
	}

	if err := server.Import(version, imageURI, filePath, connStr, logLevel); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
