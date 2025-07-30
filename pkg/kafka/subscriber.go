package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type Subscriber struct {
	r *kafka.Reader
}

func NewSubscriber(brokers []string, topic, group string) *Subscriber {
	return &Subscriber{
		r: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  group,
			MinBytes: 1e3,
			MaxBytes: 1e6,
		}),
	}
}

func (s *Subscriber) Consume(ctx context.Context, out chan<- Event) error {
	for {
		msg, err := s.r.FetchMessage(ctx)
		if err != nil {
			return err
		}
		var ev Event
		if err := json.Unmarshal(msg.Value, &ev); err == nil {
			out <- ev
		}
		if err := s.r.CommitMessages(ctx, msg); err != nil {
			return err
		}
	}
}

func (s *Subscriber) Close() error { return s.r.Close() }
