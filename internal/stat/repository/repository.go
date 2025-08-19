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
	return &StatRepository{
		Db: db,
	}
}

// func (repo *StatRepository) AddClick(linkId uint, userID uint) {
// 	today := time.Now().Format("2006-01-02")
// 	date, _ := time.Parse("2006-01-02", today)
// 	currentDate := datatypes.Date(date)

// 	var stat model.Stat
// 	err := repo.Db.First(&stat, "link_id = ? AND user_id = ? AND date = ?", linkId, userID, currentDate).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			newStat := model.Stat{
// 				LinkId: linkId,
// 				UserID: userID,
// 				Clicks: 1,
// 				Date:   currentDate,
// 			}
// 			repo.Db.Create(&newStat)
// 			return
// 		}
// 		return
// 	}

// 	stat.Clicks += 1
// 	repo.Db.Save(&stat)
// }

func (r *StatRepository) AddClick(linkID uint32, userID uint64, ts time.Time) error {
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	day := ts.UTC().Truncate(24 * time.Hour)
	currentDate := datatypes.Date(day)

	stat := model.Stat{
		LinkId: uint(linkID),
		UserID: uint(userID),
		Date:   currentDate,
		Clicks: 1,
	}

	return r.Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "link_id"}, {Name: "date"}},
		DoUpdates: clause.Assignments(map[string]any{"clicks": gorm.Expr("stats.clicks + 1")}),
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
