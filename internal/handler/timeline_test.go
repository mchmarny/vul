package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/mchmarny/vul/pkg/vul"
	"github.com/stretchr/testify/assert"
)

func TestImageTimelineHandler(t *testing.T) {
	img := url.QueryEscape("docker.io/bitnami/mariadb")
	uri := fmt.Sprintf("/api/v1/timeline?img=%s", img)

	w := testHandler(t, uri, http.MethodGet, http.StatusOK, nil)

	var out Response[map[string]*vul.ListImageTimelineItem]
	err := json.NewDecoder(w.Result().Body).Decode(&out)
	assert.NoError(t, err)
	assert.NotEmpty(t, out.Created)
	assert.NotEmpty(t, out.Criteria)
	assert.NotNil(t, out.Data)
}

func TestImageTimelineHandlerError(t *testing.T) {
	testHandler(t, "/api/v1/timeline", http.MethodGet, http.StatusBadRequest, nil)
}
