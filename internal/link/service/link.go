package service

import (
	"context"

	"github.com/sxd0/go_url-shortener/internal/link/model"
	"github.com/sxd0/go_url-shortener/internal/link/repository"
)

type LinkService struct {
	repo *repository.LinkRepository
}

func NewLinkService(repo *repository.LinkRepository) *LinkService {
	return &LinkService{
		repo: repo,
	}
}

func (s *LinkService) CreateLink(ctx context.Context, url string, userID uint) (*model.Link, error) {
	return nil, nil
}

func (s *LinkService) GetAllLinks(ctx context.Context, userID uint, limit int, offset int) ([]model.Link, error) {
	return nil, nil
}

func (s *LinkService) UpdateLink(ctx context.Context, id uint, url, hash string) (*model.Link, error) {
	return nil, nil
}

func (s *LinkService) DeleteLink(ctx context.Context, id uint) error {
	return nil
}

func (s *LinkService) GetLinkByHash(ctx context.Context, hash string) (*model.Link, error) {
	return nil, nil
}
