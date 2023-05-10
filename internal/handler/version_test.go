package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mchmarny/vul/pkg/query"
	"github.com/stretchr/testify/assert"
)

func TestImageVersionHandler(t *testing.T) {
	h := getTestHandler(t)

	vr := query.ListImageVersionRequest{
		Image: "docker.io/bitnami/mariadb",
	}

	b, err := json.Marshal(vr)
	assert.Nil(t, err)

	// request
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/versions", bytes.NewBuffer(b))
	assert.NoError(t, err)

	// execute
	h.Router.ServeHTTP(w, req)

	// validate
	assert.Equal(t, http.StatusOK, w.Code)

	var r Response[map[string][]*query.ListImageSourceItem]
	err = json.NewDecoder(w.Result().Body).Decode(&r)
	assert.NoError(t, err)
	assert.Equal(t, h.Version, r.Version)
	assert.NotEmpty(t, r.Created)
	assert.NotEmpty(t, r.Criteria)
	assert.NotNil(t, r.Data)
}
