package publisher

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	client *nats.Conn
}

func NewPublisher(url string) (*Publisher, error) {
	client, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to stream: %w", err)
	}

	return &Publisher{
		client: client,
	}, nil
}

func (p *Publisher) Publish(topic, message string) error {
	return p.client.Publish(topic, []byte(message))
}