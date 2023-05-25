package server

import (
	"fmt"
	"testing"

	"github.com/mchmarny/vul/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestImporter(t *testing.T) {
	m := map[string]string{
		"node":  "../../test/grype.json",
		"redis": "../../test/snyk.json",
		"ruby":  "../../test/trivy.json",
	}

	cnf, err := config.ReadFromFile("../../config/secret-test.yaml")
	if err != nil {
		t.Fatalf("error reading config: %v", err)
	}
	cnf.Version = "v0.0.1"

	uri := fmt.Sprintf("%s://%s:%s@/%s?host=%s",
		cnf.Store.Type, cnf.Store.User, cnf.Store.Password, cnf.Store.DB, cnf.Store.Host)

	for k, v := range m {
		t.Logf("Importing: %s - %s", k, v)
		assert.NoError(t, Import(cnf.Version, k, v, uri, cnf.Runtime.LogLevel))
	}
}
