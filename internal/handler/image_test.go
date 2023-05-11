package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageHandler(t *testing.T) {
	h := getTestHandler(t)

	// request
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/images", nil)
	assert.NoError(t, err)

	// execute
	h.Router.ServeHTTP(w, req)

	// validate
	assert.Equal(t, http.StatusOK, w.Code)

	var r Response[[]string]
	err = json.NewDecoder(w.Result().Body).Decode(&r)
	assert.NoError(t, err)
	assert.Equal(t, h.Version, r.Version)
	assert.NotEmpty(t, r.Created)
	assert.Nil(t, r.Criteria)
	assert.NotNil(t, r.Data)
}
