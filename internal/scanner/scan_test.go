package scanner

import (
	"os"
	"testing"

	"github.com/mchmarny/vul/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestAllScanners(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	dir := "test"
	assert.NoError(t, os.Mkdir(dir, 0755))
	defer os.RemoveAll(dir)

	img := "docker.io/redis@sha256:7b83a0167532d4320a87246a815a134e19e31504d85e8e55f0bb5bb9edf70448"

	cnf := config.Scanner{}

	Scan(cnf, img, dir)

	files, err := os.ReadDir(dir)
	assert.NoError(t, err)
	assert.Len(t, files, 3)
}
