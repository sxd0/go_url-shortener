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
			log.Printf("[KAFKA][CONSUME] kind=%s link_id=%d user_id=%d ts=%s",
				ev.Kind, ev.LinkID, ev.UserID, ev.Ts.UTC().Format("2006-01-02T15:04:05Z07:00"))

			if ev.Kind != kafka.LinkVisitedKind {
				continue
			}
			if err := s.repo.AddClick(uint32(ev.LinkID), uint64(ev.UserID), ev.Ts); err != nil {
				log.Printf("[STAT][UPSERT][ERR] link_id=%d user_id=%d err=%v", ev.LinkID, ev.UserID, err)
			} else {
				log.Printf("[STAT][UPSERT][OK] link_id=%d user_id=%d", ev.LinkID, ev.UserID)
			}
		}
	}
}
