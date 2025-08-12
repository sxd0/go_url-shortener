package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type Subscriber struct {
	r *kafka.Reader
}

func NewSubscriber(brokers []string, topic, group string) *Subscriber {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        group,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		MaxWait:        200 * time.Millisecond,
		CommitInterval: time.Second,
	})
	return &Subscriber{r: r}
}

func (s *Subscriber) Consume(ctx context.Context, out chan<- Event) error {
	defer close(out)

	go func() {
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				stats := s.r.Stats()
				ConsumerLag.Set(float64(stats.Lag))
			}
		}
	}()

	backoff := 100 * time.Millisecond
	for {
		m, err := s.r.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ConsumeErrors.Inc()
			time.Sleep(backoff)
			if backoff < 3*time.Second {
				backoff *= 2
			}
			continue
		}
		backoff = 100 * time.Millisecond 

		var ev Event
		if err := json.Unmarshal(m.Value, &ev); err != nil {
			ConsumeErrors.Inc()
		} else {
			ConsumeTotal.Inc()
			out <- ev
		}

		if err := s.r.CommitMessages(ctx, m); err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ConsumeErrors.Inc()
		}
	}
}

func (s *Subscriber) Close() error { return s.r.Close() }
