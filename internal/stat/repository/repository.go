package repository

import (
	"time"

	"github.com/sxd0/go_url-shortener/internal/stat/model"
	"github.com/sxd0/go_url-shortener/internal/stat/payload"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)


type StatRepository struct {
	Db *gorm.DB
}

func NewStatRepository(db *gorm.DB) *StatRepository {
	return &StatRepository{
		Db: db,
	}
}

func (repo *StatRepository) AddClick(linkId uint) {
	var stat model.Stat
	currentDate := datatypes.Date(time.Now())
	repo.Db.Find(&stat, "link_id = ? and date = ?", linkId, currentDate)
	if stat.ID == 0 {
		repo.Db.Create(&model.Stat{
			LinkId: linkId,
			Clicks: 1,
			Date:   currentDate,
		})
	} else {
		stat.Clicks += 1
		repo.Db.Save(&stat)
	}
}

func (repo *StatRepository) GetStats(by string, from, to time.Time) []payload.GetStatResponse {
	var stats []payload.GetStatResponse
	var selectQuery string
	switch by {
	case GroupByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case GroupByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	}
	repo.Db.Table("stats").
		Select(selectQuery).
		Where("date BETWEEN ? and ?", from, to).
		Group("period").
		Order("period").
		Scan(&stats)
	return stats
}

func (repo *StatRepository) GetByUserID(userID uint, from, to time.Time, groupBy string) ([]model.Stat, error) {
	var stats []model.Stat
	result := repo.Db.
		Joins("JOIN links ON stats.link_id = links.id").
		Where("links.user_id = ?", userID).
		Where("stats.created_at BETWEEN ? AND ?", from, to).
		Find(&stats)

	if result.Error != nil {
		return nil, result.Error
	}

	return stats, nil
}
