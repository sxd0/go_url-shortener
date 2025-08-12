package kafka

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	w       *kafka.Writer
	ch      chan kafka.Message
	workers int
	wg      sync.WaitGroup
	closeCh chan struct{}
}

type PubConfig struct {
	Brokers      []string
	Topic        string
	Acks         string // "0" | "1" | "all"
	BatchSize    int 
	BatchTimeout time.Duration
	Compression  string 
	QueueSize    int  
	Workers      int
}

func NewPublisherWithConfig(cfg PubConfig) *Publisher {
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 100
	}
	if cfg.BatchTimeout <= 0 {
		cfg.BatchTimeout = 100 * time.Millisecond
	}
	if cfg.QueueSize <= 0 {
		cfg.QueueSize = 1024
	}
	if cfg.Workers <= 0 {
		cfg.Workers = 2
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafka.Hash{},
		BatchSize:    cfg.BatchSize,
		BatchTimeout: cfg.BatchTimeout,
		RequiredAcks: parseAcks(cfg.Acks),
		Compression:  parseCompression(cfg.Compression),
	}

	p := &Publisher{
		w:       w,
		ch:      make(chan kafka.Message, cfg.QueueSize),
		workers: cfg.Workers,
		closeCh: make(chan struct{}),
	}

	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
	return p
}

func NewPublisher(brokers []string, topic string) *Publisher {
	return NewPublisherWithConfig(PubConfig{
		Brokers:      brokers,
		Topic:        topic,
		Acks:         "1",
		BatchSize:    200,
		BatchTimeout: 100 * time.Millisecond,
		Compression:  "snappy",
		QueueSize:    2048,
		Workers:      2,
	})
}

func (p *Publisher) Publish(_ context.Context, ev Event) error {
	data, err := json.Marshal(ev)
	if err != nil {
		PublishErrors.Inc()
		return err
	}
	msg := kafka.Message{
		Key:   buildKey(ev),
		Value: data,
		Time:  ev.Ts,
	}
	select {
	case p.ch <- msg:
		PublishTotal.Inc()
	default:
		PublishDropped.Inc()
	}
	return nil
}

func (p *Publisher) worker() {
	defer p.wg.Done()
	for {
		select {
		case <-p.closeCh:
			return
		case msg := <-p.ch:
			p.writeWithRetry(msg)
		}
	}
}

func (p *Publisher) writeWithRetry(msg kafka.Message) {
	const maxAttempts = 5
	backoff := 50 * time.Millisecond
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err := p.w.WriteMessages(ctx, msg)
		cancel()
		if err == nil {
			return
		}
		PublishErrors.Inc()
		time.Sleep(backoff)
		if backoff < 800*time.Millisecond {
			backoff *= 2
		}
	}
}

func (p *Publisher) Close() error {
	close(p.closeCh)
	close(p.ch)
	p.wg.Wait()
	return p.w.Close()
}

func buildKey(ev Event) []byte {
	if ev.LinkID > 0 {
		return []byte(strconv.FormatUint(uint64(ev.LinkID), 10))
	}
	return []byte(string(ev.Kind))
}

func parseAcks(s string) kafka.RequiredAcks {
	switch s {
	case "0":
		return kafka.RequireNone
	case "all", "-1":
		return kafka.RequireAll
	default:
		return kafka.RequireOne
	}
}

func parseCompression(s string) kafka.Compression {
	switch strings.ToLower(s) {
	case "none", "no", "off", "0":
		return 0
	case "gzip":
		return kafka.Gzip
	case "lz4":
		return kafka.Lz4
	case "zstd":
		return kafka.Zstd
	case "snappy", "":
		fallthrough
	default:
		return kafka.Snappy
	}
}

