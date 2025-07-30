package service

import (
	"context"
	"errors"
	"time"

	"github.com/sxd0/go_url-shortener/internal/link/model"
	"github.com/sxd0/go_url-shortener/internal/link/repository"
	"github.com/sxd0/go_url-shortener/pkg/kafka"
)

const (
	maxHashGenerateAttempts = 5
	hashLength              = 10
)

type LinkService struct {
	repo      *repository.LinkRepository
	publisher *kafka.Publisher
}

func NewLinkService(repo *repository.LinkRepository, pub *kafka.Publisher) *LinkService {
	return &LinkService{repo: repo, publisher: pub}
}

func (s *LinkService) CreateLink(ctx context.Context, url string, userID uint) (*model.Link, error) {
	for i := 0; i < maxHashGenerateAttempts; i++ {
		link, err := model.NewLink(url, func(hash string) bool {
			exists, _ := s.repo.GetByHash(hash)
			return exists != nil
		})
		if err != nil {
			return nil, err
		}
		link.UserID = userID

		created, err := s.repo.Create(link)
		if err == nil {
			if s.publisher != nil {
				_ = s.publisher.Publish(ctx, kafka.Event{
					Kind:   kafka.LinkCreatedKind,
					LinkID: created.ID,
					UserID: userID,
					Ts:     time.Now().UTC(),
				})
			}
			return created, nil
		}
	}
	return nil, errors.New("failed to create unique link after several attempts")
}

func (s *LinkService) GetAllLinks(ctx context.Context, userID uint, limit int, offset int) ([]model.Link, int64, error) {
	links, err := s.repo.GetAllByUserID(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.CountByUserID(userID)
	if err != nil {
		return nil, 0, err
	}
	return links, count, nil
}

func (s *LinkService) UpdateLink(ctx context.Context, id uint, url, hash string) (*model.Link, error) {
	var link *model.Link
	var err error

	if id != 0 {
		link, err = s.repo.GetByID(id)
	} else if hash != "" {
		link, err = s.repo.GetByHash(hash)
	} else {
		return nil, errors.New("either id or hash is required")
	}
	if err != nil || link == nil {
		return nil, errors.New("link not found")
	}

	link.Url = url
	if hash != "" {
		link.Hash = hash
	}
	return s.repo.Update(link)
}

func (s *LinkService) DeleteLink(ctx context.Context, id uint) error {
	return s.repo.Delete(id)
}

func (s *LinkService) GetLinkByHash(ctx context.Context, hash string) (*model.Link, error) {
	return s.repo.GetByHash(hash)
}
