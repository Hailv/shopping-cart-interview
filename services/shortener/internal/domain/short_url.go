package domain

import (
	"errors"
	"regexp"
)

type ShortURL struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	OriginalURL string `gorm:"type:text;not null"`
	URLHash     string `gorm:"type:char(64);not null"`
	ShortCode   string `gorm:"type:varchar(10);unique;not null"`
}

func NewShortURL(originalURL, urlHash string) (*ShortURL, error) {
	if !isValidURL(originalURL) {
		return nil, errors.New("Invalid url")
	}

	return &ShortURL{
		OriginalURL: originalURL,
		URLHash:     urlHash,
	}, nil
}

func isValidURL(url string) bool {
	var urlRegex = regexp.MustCompile(
		`^(https?):\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$`,
	)
	return urlRegex.MatchString(url)
}
