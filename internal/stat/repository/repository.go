package repository

import (
	"errors"
	"log"
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
	today := time.Now().Format("2006-01-02") 
	date, err := time.Parse("2006-01-02", today)
	if err != nil {
		log.Println("Failed to parse date:", err)
		return
	}

	currentDate := datatypes.Date(date)

	var stat model.Stat
	err = repo.Db.First(&stat, "link_id = ? AND date = ?", linkId, currentDate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newStat := model.Stat{
				LinkId: linkId,
				Clicks: 1,
				Date:   currentDate,
			}
			if err := repo.Db.Create(&newStat).Error; err != nil {
				log.Println("Failed to create stat:", err)
			}
			return
		}
		log.Println("Failed to query stat:", err)
		return
	}

	stat.Clicks += 1
	if err := repo.Db.Save(&stat).Error; err != nil {
		log.Println("Failed to update stat:", err)
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

