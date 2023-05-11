package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mchmarny/vul/pkg/vul"
	"github.com/stretchr/testify/assert"
)

func TestImageSummaryHandler(t *testing.T) {
	h := getTestHandler(t)

	validateImageSummaryResponse(t, h, "")
	validateImageSummaryResponse(t, h, "docker.io/bitnami/mariadb")
}

func validateImageSummaryResponse(t *testing.T, h *Handler, uri string) {
	var body io.Reader
	method := http.MethodGet

	if uri != "" {
		b, err := json.Marshal(vul.ImageRequest{
			Image: uri,
		})
		assert.Nil(t, err)
		body = bytes.NewReader(b)
		method = http.MethodPost
	}

	req, err := http.NewRequest(method, "/api/v1/summary", body)
	assert.NoError(t, err)

	// execute
	w := httptest.NewRecorder()
	h.Router.ServeHTTP(w, req)

	// validate
	assert.Equal(t, http.StatusOK, w.Code)

	var r Response[*vul.ImageRequest]
	err = json.NewDecoder(w.Result().Body).Decode(&r)
	assert.NoError(t, err)
	assert.Equal(t, h.Version, r.Version)
	assert.NotEmpty(t, r.Created)
	assert.Nil(t, r.Criteria)
	assert.NotNil(t, r.Data)
}
