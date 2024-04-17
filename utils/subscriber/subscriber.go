package subscriber

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

var timeout = 30 * time.Minute

type Subscriber struct {
	sub *nats.Subscription
}

func NewSubscriber(url string, topic string) (*Subscriber, error) {
	client, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to stream: %w", err)
	}

	sub, err := client.SubscribeSync(topic)
	if err != nil {
		return nil, fmt.Errorf("couldn't subscribe on topic %s: %w", topic, err)
	}
	

	return &Subscriber{
		sub: sub,
	}, nil
}

func (s *Subscriber) GetMessage() (string, error) {
	msg, err := s.sub.NextMsg(timeout)

	return string(msg.Data), err
}
