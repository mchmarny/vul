package handler

import (
	"encoding/json"
	"net/http"
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
	var in interface{}
	method := http.MethodGet

	if uri != "" {
		in = vul.ImageRequest{
			Image: uri,
		}
		method = http.MethodPost
	}

	w := testHandler(t, "/api/v1/summary", method, http.StatusOK, in)

	var out Response[*vul.ImageRequest]
	err := json.NewDecoder(w.Result().Body).Decode(&out)
	assert.NoError(t, err)
	assert.NotEmpty(t, out.Created)
	assert.Nil(t, out.Criteria)
	assert.NotNil(t, out.Data)
}
