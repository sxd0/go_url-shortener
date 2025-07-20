package link

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/sxd0/go_url-shortener/internal/auth/model"
	"github.com/sxd0/go_url-shortener/internal/stat"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Url    string `json:"url"`
	Hash   string `json:"hash" gorm:"uniqueIndex"`
	UserID uint
	User   model.User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Stats  []stat.Stat `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

const hashLength = 10
const hashAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const maxHashGenerateAttempts = 5

func NewLink(url string, checkExists func(hash string) bool) (*Link, error) {
	for i := 0; i < maxHashGenerateAttempts; i++ {
		hash, err := generateSecureHash()
		if err != nil {
			return nil, err
		}
		if !checkExists(hash) {
			return &Link{
				Url:  url,
				Hash: hash,
			}, nil
		}
	}
	return nil, errors.New("failed to generate unique hash after several attempts")
}

func generateSecureHash() (string, error) {
	b := make([]byte, hashLength)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(hashAlphabet))))
		if err != nil {
			return "", err
		}
		b[i] = hashAlphabet[num.Int64()]
	}
	return string(b), nil
}
