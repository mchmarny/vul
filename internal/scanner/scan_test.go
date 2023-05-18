package scanner

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllScanners(t *testing.T) {
	dir := "test"
	assert.NoError(t, os.Mkdir(dir, 0755))
	defer os.RemoveAll(dir)

	img := "docker.io/redis@sha256:7b83a0167532d4320a87246a815a134e19e31504d85e8e55f0bb5bb9edf70448"

	Scan(img, dir)

	files, err := os.ReadDir(dir)
	assert.NoError(t, err)
	assert.Len(t, files, 3)
}
