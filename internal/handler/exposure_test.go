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

func TestImageVersionExposureHandler(t *testing.T) {
	img := url.QueryEscape("docker.io/bitnami/mariadb")
	dig := url.QueryEscape("sha256:97b0be98b4714e81dac9ac55513f4f87c627d88da09d90c708229835124a8215")
	uri := fmt.Sprintf("/api/v1/exposures?img=%s&dig=%s", img, dig)

	w := testHandler(t, uri, http.MethodGet, http.StatusOK, nil)

	var out Response[map[string][]*vul.ListDigestExposureItem]
	err := json.NewDecoder(w.Result().Body).Decode(&out)
	assert.NoError(t, err)
	assert.NotEmpty(t, out.Created)
	assert.NotEmpty(t, out.Criteria)
	assert.NotEmpty(t, out.Data)
}
