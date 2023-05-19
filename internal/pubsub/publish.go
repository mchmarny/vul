package pubsub

import (
	"context"
	"encoding/json"

	api "cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Publisher interface {
	Publish(ctx context.Context, topic string, content interface{}) error
	PublishStr(ctx context.Context, topic, content string) error
}

// New creates a new PubSub publisher.
func New(ctx context.Context, projectID string, send bool) (Publisher, error) {
	if projectID == "" {
		return nil, errors.New("conf is nil or project ID string is empty")
	}

	s := &SimplePublisher{
		projectID: projectID,
		topics:    make(map[string]*api.Topic),
	}

	if send {
		client, err := api.NewClient(ctx, projectID)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating PubSub client for project: %s", projectID)
		}
		s.client = client
	}

	log.Debug().Msgf("pubsub publisher created in %s project", projectID)

	return s, nil
}

type SimplePublisher struct {
	projectID string
	client    *api.Client
	topics    map[string]*api.Topic
}

func (p *SimplePublisher) Publish(ctx context.Context, topic string, content interface{}) error {
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

	if err := p.publish(ctx, topic, b); err != nil {
		return errors.Wrapf(err, "error publishing content to: %s", topic)
	}

	return nil
}

func (p *SimplePublisher) PublishStr(ctx context.Context, topic, content string) error {
	if topic == "" {
		return errors.New("topic is empty")
	}
	if content == "" {
		return errors.New("content is nil")
	}

	if err := p.publish(ctx, topic, []byte(content)); err != nil {
		return errors.Wrapf(err, "error publishing string content to: %s", topic)
	}

	return nil
}

func (p *SimplePublisher) publish(ctx context.Context, topic string, content []byte) error {
	if topic == "" {
		return errors.New("topic is empty")
	}
	if content == nil {
		return errors.New("content is nil")
	}

	if p.client == nil {
		log.Debug().
			Str("topic", topic).
			Str("content", string(content)).
			Msgf("log only, send is disabled")
		return nil
	}

	t, ok := p.topics[topic]
	if !ok {
		t = p.client.Topic(topic)
		p.topics[topic] = t
	}

	result := t.Publish(ctx, &api.Message{Data: content})
	id, err := result.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting published message ID")
	}

	log.Debug().Msgf("message published, ID: %s", id)

	return nil
}
