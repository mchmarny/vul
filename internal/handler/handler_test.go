package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestHandler(t *testing.T) *Handler {
	h, err := New(context.Background(), "test", "v0.0.1")
	assert.NoError(t, err)
	assert.NotNil(t, h)
	return h
}

func TestHandler(t *testing.T) {
	h := getTestHandler(t)

	// request
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	// execute
	h.Router.ServeHTTP(w, req)

	// validate
	assert.Equal(t, http.StatusOK, w.Code)
}
