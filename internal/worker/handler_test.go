package worker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mchmarny/vul/internal/config"
	"github.com/stretchr/testify/assert"
)

func getTestHandler(t *testing.T) *Handler {
	cnf, err := config.ReadFromFile("../../config/secret-test.yaml")
	if err != nil {
		t.Fatalf("error reading config: %v", err)
	}
	cnf.Version = "v0.0.1"

	h, err := New(context.Background(), cnf)
	assert.NoError(t, err)
	assert.NotNil(t, h)
	return h
}

func TestHealthHandler(t *testing.T) {
	testHandler(t, "/health", http.StatusOK)
}

func testHandler(t *testing.T, url string, expectedCode int) *httptest.ResponseRecorder {
	h := getTestHandler(t)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)
	h.Router.ServeHTTP(w, req)
	assert.Equal(t, expectedCode, w.Code)
	return w
}
