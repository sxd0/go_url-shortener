package service

import (
	"context"
	"log"

	"github.com/sxd0/go_url-shortener/internal/stat/repository"
	"github.com/sxd0/go_url-shortener/pkg/kafka"
)

type StatService struct {
	repo *repository.StatRepository
	sub  *kafka.Subscriber
}

func NewStatService(repo *repository.StatRepository, sub *kafka.Subscriber) *StatService {
	return &StatService{repo: repo, sub: sub}
}

func (s *StatService) Start(ctx context.Context) error {
	events := make(chan kafka.Event, 128)
	go s.sub.Consume(ctx, events)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev := <-events:
			if ev.Kind != kafka.LinkVisitedKind {
				continue
			}
			if err := s.repo.AddClick(uint32(ev.LinkID), uint64(ev.UserID)); err != nil {
				log.Printf("stat upsert: %v", err)
			}
		}
	}
}
