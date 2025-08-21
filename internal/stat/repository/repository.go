package repository

import (
	"time"

	"github.com/sxd0/go_url-shortener/internal/stat/model"
	"github.com/sxd0/go_url-shortener/internal/stat/payload"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

type StatRepository struct {
	Db *gorm.DB
}

func NewStatRepository(db *gorm.DB) *StatRepository {
	return &StatRepository{Db: db}
}

func (r *StatRepository) AddClick(linkID uint32, userID uint64, ts time.Time) error {
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	day := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, time.UTC)
	currentDate := datatypes.Date(day)

	stat := model.Stat{
		LinkId: uint(linkID),
		UserID: uint(userID),
		Date:   currentDate,
		Clicks: 1,
	}

	return r.Db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "link_id"}, {Name: "date"}},
		DoUpdates: clause.Assignments(map[string]any{
			"clicks":     gorm.Expr("stats.clicks + EXCLUDED.clicks"),
			"updated_at": gorm.Expr("NOW()"),
		}),
	}).Create(&stat).Error
}

func (r *StatRepository) GetStats(userID uint, by string, from, to time.Time) []payload.GetStatResponse {
	var stats []payload.GetStatResponse
	var selectQuery string
	switch by {
	case GroupByDay:
		selectQuery = `link_id, to_char(date, 'YYYY-MM-DD') AS date, sum(clicks) AS clicks`
	case GroupByMonth:
		selectQuery = `link_id, to_char(date, 'YYYY-MM') AS date, sum(clicks) AS clicks`
	default:
		selectQuery = `link_id, to_char(date, 'YYYY-MM-DD') AS date, sum(clicks) AS clicks`
	}

	r.Db.Table("stats").
		Select(selectQuery).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, from, to).
		Group("link_id, date").
		Order("link_id, date").
		Scan(&stats)

	return stats
}
