package service

import (
	"log"

	"github.com/sxd0/go_url-shortener/internal/stat/repository"
	eventPayload "github.com/sxd0/go_url-shortener/internal/stat/event"
	"github.com/sxd0/go_url-shortener/pkg/event"
)

type StatServiceDeps struct {
	EventBus       *event.EventBus
	StatRepository *repository.StatRepository
}

type StatService struct {
	EventBus       *event.EventBus
	StatRepository *repository.StatRepository
}

func NewStatService(deps *StatServiceDeps) *StatService {
	return &StatService{
		EventBus:       deps.EventBus,
		StatRepository: deps.StatRepository,
	}
}

func (s *StatService) AddClick() {
	for msg := range s.EventBus.Subscribe() {
		if msg.Type != event.EventLinkVisited {
			continue
		}

		payload, ok := msg.Data.(eventPayload.LinkVisitedEvent)
		if !ok {
			log.Printf("EventLinkVisited: unexpected data type: %T\n", msg.Data)
			continue
		}

		s.StatRepository.AddClick(payload.LinkID, payload.UserID)
	}
}
