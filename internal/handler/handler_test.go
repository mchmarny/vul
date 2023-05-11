package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
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
	testHandler(t, "/health", http.MethodGet, http.StatusOK, nil)
}

func testHandler(t *testing.T, url, method string, expectedCode int, d interface{}) *httptest.ResponseRecorder {
	h := getTestHandler(t)
	var r io.Reader

	if d != nil {

		b, err := json.Marshal(d)
		assert.Nil(t, err)
		r = bytes.NewBuffer(b)
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, url, r)
	assert.NoError(t, err)
	h.Router.ServeHTTP(w, req)
	assert.Equal(t, expectedCode, w.Code)
	return w
}
