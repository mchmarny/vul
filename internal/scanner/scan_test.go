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

	img := "docker.io/redis"

	cnf := config.Scanner{}

	Scan(cnf, img, dir)

	files, err := os.ReadDir(dir)
	assert.NoError(t, err)
	assert.Len(t, files, 3)
}
