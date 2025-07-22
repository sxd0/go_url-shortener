package repository

import (
	"github.com/sxd0/go_url-shortener/internal/link/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// type LinkRepositoryDeps struct {
// 	DataBase *db.Db
// }

type LinkRepository struct {
	Database *gorm.DB
}

func NewLinkRepository(database *gorm.DB) *LinkRepository {
	return &LinkRepository{
		Database: database,
	}
}

func (repo *LinkRepository) Create(link *model.Link) (*model.Link, error) {
	result := repo.Database.Create(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repo *LinkRepository) GetByHash(hash string) (*model.Link, error) {
	var link model.Link
	result := repo.Database.Where("hash = ?", hash).First(&link)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (repo *LinkRepository) Update(link *model.Link) (*model.Link, error) {
	result := repo.Database.Clauses(clause.Returning{}).Updates(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repo *LinkRepository) Delete(id uint) error {
	result := repo.Database.Delete(&model.Link{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *LinkRepository) GetById(id uint) (*model.Link, error) {
	var link model.Link
	result := repo.Database.First(&link, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (repo *LinkRepository) Count() int64 {
	var count int64
	repo.Database.
		Table("links").
		Where("deleted_at is null").
		Count(&count)
	return count
}

func (repo *LinkRepository) GetAll(limit, offset int) []model.Link {
	var links []model.Link

	repo.Database.
		Table("links").
		Where("deleted_at is null").
		Order("id asc").
		Limit(limit).
		Offset(offset).
		Scan(&links)
	return links
}

func (repo *LinkRepository) GetAllByUserID(userID uint, limit int, offset int) ([]model.Link, error) {
	var links []model.Link
	result := repo.Database.Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&links)
	if result.Error != nil {
		return nil, result.Error
	}
	return links, nil
}
