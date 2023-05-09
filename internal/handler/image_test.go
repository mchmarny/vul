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

	return &Handler{
		Name:    "test",
		Version: "1.0.0",
		Pool:    pool,
	}
}

func TestImageHandler(t *testing.T) {
	h := getTestHandler(t)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ImageHandler)

	handler.ServeHTTP(rr, req)
	status := rr.Code
	assert.Equal(t, http.StatusOK, status)

	var r Response[[]*query.ListImageItem]
	err = json.NewDecoder(rr.Body).Decode(&r)
	assert.NoError(t, err)
	assert.Equal(t, h.Version, r.Version)
	assert.NotEmpty(t, r.Created)
	assert.NotNil(t, r.Data)
}
