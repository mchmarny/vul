package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c, err := ReadFromFile("../../config/secret-test.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.NotEmpty(t, c.Name)
	assert.NotEmpty(t, c.ProjectID)
	assert.NotNil(t, c.Runtime)
	assert.NotEmpty(t, c.Runtime.LogLevel)
	assert.NotEmpty(t, c.Runtime.ExternalURL)
	assert.NotNil(t, c.Store)
	assert.NotEmpty(t, c.Store.URI)
}
