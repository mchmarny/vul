package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mchmarny/vul/pkg/vul"
	"github.com/stretchr/testify/assert"
)

func TestImageTimelineHandler(t *testing.T) {
	in := vul.ListImageTimelineRequest{
		Image: "docker.io/bitnami/mongodb",
	}
	w := testHandler(t, "/api/v1/timeline", http.MethodPost, http.StatusOK, in)

	var out Response[map[string]*vul.ListImageTimelineItem]
	err := json.NewDecoder(w.Result().Body).Decode(&out)
	assert.NoError(t, err)
	assert.NotEmpty(t, out.Created)
	assert.NotEmpty(t, out.Criteria)
	assert.NotNil(t, out.Data)
}

func TestImageTimelineHandlerError(t *testing.T) {
	testHandler(t, "/api/v1/timeline", http.MethodPost, http.StatusBadRequest, nil)
}

func TestImageTimelineHandlerEmptyError(t *testing.T) {
	r := vul.ListImageTimelineRequest{}
	testHandler(t, "/api/v1/timeline", http.MethodPost, http.StatusBadRequest, r)
}
