package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/pkg/query"
	"github.com/stretchr/testify/assert"
)

func getTestHandler(t *testing.T) *Handler {
	pool, err := data.GetPool(context.Background(), os.Getenv("DATA_URI"))
	if err != nil {
		t.Fatalf("error getting data pool: %v", err)
	}

	return New("test", "v0.0.1", pool)
}

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

	var r Response[[]*query.ListImageItem]
	err = json.NewDecoder(w.Result().Body).Decode(&r)
	assert.NoError(t, err)
	assert.Equal(t, h.Version, r.Version)
	assert.NotEmpty(t, r.Created)
	assert.Nil(t, r.Criteria)
	assert.NotNil(t, r.Data)
}
