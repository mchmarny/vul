package pubsub

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"google.golang.org/api/pubsub/v1"
)

// PubSubMessage is the payload of a Pub/Sub event.
type pushRequest struct {
	Message      pubsub.PubsubMessage `json:"message"`
	Subscription string               `json:"subscription"`
}

func GetData(c *gin.Context) ([]byte, error) {
	var t pushRequest
	if err := c.ShouldBindJSON(&t); err != nil {
		return nil, errors.Wrap(err, "invalid task format")
	}

	if t.Message.Data == "" {
		return nil, errors.New("invalid request, empty data field")
	}

	b, err := base64.StdEncoding.DecodeString(t.Message.Data)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding message data")
	}

	return b, nil
}

func BindMessage[T any](c *gin.Context) (*T, error) {
	b, err := GetData(c)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding message data")
	}

	var r T
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling message data")
	}

	return &r, nil
}
