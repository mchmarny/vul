package server

import (
	"testing"
)

func TestHealthHandler(t *testing.T) {
	RunImport("v0.0.1", "redis", "")
}
