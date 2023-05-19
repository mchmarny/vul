package processor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mchmarny/vul/internal/config"
	"github.com/mchmarny/vul/internal/pubsub"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func getTestHandler(t *testing.T) *Handler {
	cnf, err := config.ReadFromFile("../../config/secret-test.yaml")
	if err != nil {
		t.Fatalf("error reading config: %v", err)
	}
	cnf.Version = "v0.0.1"

	pubsubNew = testPublisherProvider

	h, err := New(context.Background(), cnf)
	assert.NoError(t, err)
	assert.NotNil(t, h)
	return h
}

func TestHealthHandler(t *testing.T) {
	testHandler(t, "/health", http.StatusOK)
}

func testHandler(t *testing.T, url string, expectedCode int) *httptest.ResponseRecorder {
	h := getTestHandler(t)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)
	h.Router.ServeHTTP(w, req)
	assert.Equal(t, expectedCode, w.Code)
	return w
}

func testPublisherProvider(_ context.Context, projectID string) (pubsub.Publisher, error) {
	p := &testPublisher{
		projectID: projectID,
	}
	return p, nil
}

type testPublisher struct {
	projectID string
}

func (p *testPublisher) Publish(_ context.Context, topic string, content interface{}) error {
	if topic == "" {
		return errors.New("topic is empty")
	}
	if content == nil {
		return errors.New("content is nil")
	}

	b, err := json.Marshal(content)
	if err != nil {
		return errors.Wrapf(err, "error marshaling content: %+v", content)
	}

	log.Debug().Str("topic", topic).Str("content", string(b)).Msg("publishing")

	return nil
}

func (p *testPublisher) PublishStr(_ context.Context, topic, content string) error {
	if topic == "" {
		return errors.New("topic is empty")
	}
	if content == "" {
		return errors.New("content is nil")
	}

	log.Debug().Str("topic", topic).Str("content", content).Msg("publishing")

	return nil
}
