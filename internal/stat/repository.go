package stat

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
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
	var stat Stat
	currentDate := datatypes.Date(time.Now())
	repo.Db.Find(&stat, "link_id = ? and date = ?", linkId, currentDate)
	if stat.ID == 0 {
		repo.Db.Create(&Stat{
			LinkId: linkId,
			Clicks: 1,
			Date:   currentDate,
		})
	} else {
		stat.Clicks += 1
		repo.Db.Save(&stat)
	}
}

func (repo *StatRepository) GetStats(by string, from, to time.Time) []GetStatResponse {
	var stats []GetStatResponse
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

func (repo *StatRepository) GetByUserID(userID uint, from, to time.Time, groupBy string) ([]Stat, error) {
	var stats []Stat
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
