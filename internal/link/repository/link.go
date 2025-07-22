package repository

import (
	"github.com/sxd0/go_url-shortener/internal/link/model"
	"gorm.io/gorm"
)

type LinkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) *LinkRepository {
	return &LinkRepository{
		db: db,
	}
}

func (r *LinkRepository) Create(link *model.Link) (*model.Link, error) {
	if err := r.db.Create(link).Error; err != nil {
		return nil, err
	}
	return link, nil
}

func (r *LinkRepository) GetByHash(hash string) (*model.Link, error) {
	var link model.Link
	if err := r.db.Where("hash = ?", hash).First(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepository) Update(link *model.Link) (*model.Link, error) {
	if err := r.db.Save(link).Error; err != nil {
		return nil, err
	}
	return link, nil
}

func (r *LinkRepository) Delete(id uint) error {
	return r.db.Delete(&model.Link{}, id).Error
}

func (r *LinkRepository) GetByID(id uint) (*model.Link, error) {
	var link model.Link
	if err := r.db.First(&link, id).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepository) GetAllByUserID(userID uint, limit int, offset int) ([]model.Link, error) {
	var links []model.Link
	if err := r.db.Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

func (r *LinkRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Link{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
