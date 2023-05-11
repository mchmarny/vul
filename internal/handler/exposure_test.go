package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mchmarny/vul/pkg/vul"
	"github.com/stretchr/testify/assert"
)

func TestImageVersionExposureHandler(t *testing.T) {
	h := getTestHandler(t)

	vr := vul.ListImageVersionExposureRequest{
		Image:  "docker.io/bitnami/mongodb",
		Digest: "sha256:419f129df0140834d89c94b29700c91f38407182137be480a0d6c6cbe2e0d00a",
	}

	b, err := json.Marshal(vr)
	assert.Nil(t, err)

	// request
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/exposures", bytes.NewBuffer(b))
	assert.NoError(t, err)

	// execute
	h.Router.ServeHTTP(w, req)

	// validate
	assert.Equal(t, http.StatusOK, w.Code)

	var r Response[map[string][]*vul.ListDigestExposureItem]
	err = json.NewDecoder(w.Result().Body).Decode(&r)
	assert.NoError(t, err)
	assert.Equal(t, h.Version, r.Version)
	assert.NotEmpty(t, r.Created)
	assert.NotEmpty(t, r.Criteria)
	assert.NotNil(t, r.Data)
}
