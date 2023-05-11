package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mchmarny/vul/pkg/vul"
	"github.com/stretchr/testify/assert"
)

func TestImageVersionHandler(t *testing.T) {
	in := vul.ImageRequest{
		Image: "docker.io/bitnami/mariadb",
	}

	w := testHandler(t, "/api/v1/versions", http.MethodPost, http.StatusOK, in)

	var out Response[map[string][]*vul.ListImageSourceItem]
	err := json.NewDecoder(w.Result().Body).Decode(&out)
	assert.NoError(t, err)
	assert.NotEmpty(t, out.Created)
	assert.NotEmpty(t, out.Criteria)
	assert.NotNil(t, out.Data)
}

func TestImageVersionHandlerError(t *testing.T) {
	testHandler(t, "/api/v1/versions", http.MethodPost, http.StatusBadRequest, nil)
}

func TestImageVersionHandlerEmpoty(t *testing.T) {
	r := vul.ImageRequest{}
	testHandler(t, "/api/v1/versions", http.MethodPost, http.StatusBadRequest, r)
}
