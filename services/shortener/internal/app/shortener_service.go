package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/cinchprotocol/cinch-api/services/shortener/internal/domain"
	"github.com/cinchprotocol/cinch-api/services/shortener/internal/infrastructure/id_generator"
)

var ErrURLNotFound = errors.New("URL Not found")

type Hasher interface {
	Hash(input string) string
}

type SHA256Hasher struct{}

func (s *SHA256Hasher) Hash(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// ShortenerService is the main logic
type ShortenerService struct {
	repo   domain.ShortURLRepository
	hasher Hasher
}

func NewShortenerService(repo domain.ShortURLRepository, hasher Hasher) *ShortenerService {
	return &ShortenerService{
		repo:   repo,
		hasher: hasher,
	}
}

// ShortenURL reduce the url length and save it
func (s *ShortenerService) ShortenURL(ctx context.Context, originalURL string) (string, error) {
	urlHash := s.hasher.Hash(originalURL)

	// check existing hashed
	existingURLsWithHash, err := s.repo.FindByURLHash(ctx, urlHash)
	if err != nil {
		return "", errors.New("lỗi khi kiểm tra URL hash: " + err.Error())
	}

	foundExactMatch := false
	var existingShortCode string

	if len(existingURLsWithHash) > 0 {
		// incase found any hashed
		for _, urlEntry := range existingURLsWithHash {
			if urlEntry.OriginalURL == originalURL {
				existingShortCode = urlEntry.ShortCode
				foundExactMatch = true
				break
			}
		}
	}

	if foundExactMatch {
		return existingShortCode, nil
	}

	// save if no hashed found or no collision
	newShortURL, err := domain.NewShortURL(originalURL, urlHash)
	if err != nil {
		return "", err
	}

	generatedID, err := s.repo.Insert(ctx, newShortURL)
	if err != nil {
		return "", errors.New("Error when saving data: " + err.Error())
	}

	shortCode := id_generator.Base62Encode(generatedID)

	err = s.repo.UpdateShortCode(ctx, generatedID, shortCode)
	if err != nil {
		return "", errors.New("Error when saving shortCode: " + err.Error())
	}

	return shortCode, nil
}

// GetOriginalURL get url via shortcode
func (s *ShortenerService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	shortURL, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}
	if shortURL == nil {
		return "", ErrURLNotFound
	}
	return shortURL.OriginalURL, nil
}

// GetAllShortURLs Just a test function
func (s *ShortenerService) GetAllShortURLs(ctx context.Context) ([]*domain.ShortURL, error) {
	return s.repo.FindAll(ctx)
}
