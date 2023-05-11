package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/mchmarny/vul/pkg/vul"
	"github.com/stretchr/testify/assert"
)

func TestImageSummaryHandler(t *testing.T) {
	h := getTestHandler(t)

	validateImageSummaryResponse(t, h, "")
	validateImageSummaryResponse(t, h, "docker.io/bitnami/mariadb")
}

func validateImageSummaryResponse(t *testing.T, h *Handler, img string) {
	uri := "/api/v1/summary"

	if img != "" {
		uri += "?img=" + url.QueryEscape(img)
	}

	w := testHandler(t, uri, http.MethodGet, http.StatusOK, nil)

	var out Response[*vul.SummaryItem]
	err := json.NewDecoder(w.Result().Body).Decode(&out)
	assert.NoError(t, err)
	assert.NotEmpty(t, out.Created)
	assert.Nil(t, out.Criteria)
	assert.NotNil(t, out.Data)
	assert.Equal(t, img, out.Data.Image)
}
