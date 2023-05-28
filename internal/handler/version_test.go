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

const (
	testImage       = "docker.io/bitnami/mariadb"
	imageVersionURL = "/api/v1/versions"
)

func TestImageVersionHandler(t *testing.T) {
	uri := fmt.Sprintf("%s?img=%s",
		imageVersionURL, url.QueryEscape("docker.io/bitnami/mariadb"))

	w := testHandler(t, uri, http.StatusOK)

	var out Response[[]*vul.ImageVersion]
	err := json.NewDecoder(w.Result().Body).Decode(&out)
	assert.NoError(t, err)
	assert.NotEmpty(t, out.Created)
	assert.NotEmpty(t, out.Criteria)
	assert.NotEmpty(t, out.Data)
}

func TestImageVersionHandlerError(t *testing.T) {
	testHandler(t, imageVersionURL, http.StatusBadRequest)
}
