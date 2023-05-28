package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageHandler(t *testing.T) {
	w := testHandler(t, "/api/v1/images", http.StatusOK)

	var out Response[[]string]
	err := json.NewDecoder(w.Result().Body).Decode(&out)
	assert.NoError(t, err)
	assert.NotEmpty(t, out.Created)
	assert.Nil(t, out.Criteria)
	assert.NotEmpty(t, out.Data)
}
