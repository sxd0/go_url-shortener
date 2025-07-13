package stat

import (
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

type Stat struct {
	gorm.Model
	LinkId uint           `json:"link_id"`
	Clicks int            `json:"clicks"`
	Date   datatypes.Date `json:"date"`
}
