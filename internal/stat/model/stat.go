package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Stat struct {
	gorm.Model
	LinkId uint           `json:"link_id" gorm:"not null;index;uniqueIndex:idx_link_date"`
	UserID uint           `json:"user_id" gorm:"not null"`
	Clicks int            `json:"clicks" gorm:"not null"`
	Date   datatypes.Date `json:"date" gorm:"not null;index;uniqueIndex:idx_link_date"`
}
