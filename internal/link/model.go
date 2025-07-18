package link

import (
	"go/test-http/internal/stat"
	"math/rand"
	"go/test-http/internal/user"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Url    string `json:"url"`
	Hash   string `json:"hash" gorm:"uniqueIndex"`
	UserID uint
	User   user.User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Stats  []stat.Stat `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func NewLink(url string) *Link {
	link := &Link{
		Url: url,
	}
	link.GenerateHash()
	return link
}

func (link *Link) GenerateHash() {
	link.Hash = RangStringRunes(10)
}

var letterRunes = []rune("abcdefghijklmnoprstuvwxyzABCDEFGHIGKLMNOPRSTUVWXYZ")

func RangStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
