package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mchmarny/vul/pkg/vul"
	"github.com/stretchr/testify/assert"
)

func TestImageVersionExposureHandler(t *testing.T) {
	in := vul.ListImageVersionExposureRequest{
		Image:  "docker.io/bitnami/mongodb",
		Digest: "sha256:419f129df0140834d89c94b29700c91f38407182137be480a0d6c6cbe2e0d00a",
	}

	w := testHandler(t, "/api/v1/exposures", http.MethodPost, http.StatusOK, in)

	var out Response[map[string][]*vul.ListDigestExposureItem]
	err := json.NewDecoder(w.Result().Body).Decode(&out)
	assert.NoError(t, err)
	assert.NotEmpty(t, out.Created)
	assert.NotEmpty(t, out.Criteria)
	assert.NotNil(t, out.Data)
}
