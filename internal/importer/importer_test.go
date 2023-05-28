package importer

import (
	"context"
	"testing"

	"github.com/mchmarny/vul/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestImporter(t *testing.T) {
	cnf, err := config.ReadFromFile("../../config/secret-test.yaml")
	if err != nil {
		t.Fatalf("error reading config: %v", err)
	}
	cnf.Version = "v0.0.1"

	ctx := context.Background()

	m := map[string]string{
		"node":  "../../test/grype.json",
		"redis": "../../test/snyk.json",
		"ruby":  "../../test/trivy.json",
	}

	for k, v := range m {
		t.Logf("Importing: %s - %s", k, v)

		opt, err := ParseOptions(ctx, cnf, k, v)
		assert.NoError(t, err)

		err = Import(ctx, opt)
		assert.NoError(t, err)
	}
}
