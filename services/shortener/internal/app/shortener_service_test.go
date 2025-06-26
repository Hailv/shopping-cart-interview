package app_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/cinchprotocol/cinch-api/services/shortener/internal/app"
	"github.com/cinchprotocol/cinch-api/services/shortener/internal/domain"
)

type MockShortURLRepository struct {
	mock.Mock
}

func (m *MockShortURLRepository) Insert(ctx context.Context, shortURL *domain.ShortURL) (uint64, error) {
	args := m.Called(ctx, shortURL)
	id := args.Get(0).(uint64)
	return id, args.Error(1)
}

func (m *MockShortURLRepository) UpdateShortCode(ctx context.Context, id uint64, shortCode string) error {
	args := m.Called(ctx, id, shortCode)
	return args.Error(0)
}

func (m *MockShortURLRepository) FindByShortCode(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	args := m.Called(ctx, shortCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ShortURL), args.Error(1)
}

func (m *MockShortURLRepository) FindByURLHash(ctx context.Context, urlHash string) ([]*domain.ShortURL, error) {
	args := m.Called(ctx, urlHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ShortURL), args.Error(1)
}

func (m *MockShortURLRepository) FindAll(ctx context.Context) ([]*domain.ShortURL, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ShortURL), args.Error(1)
}

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) Hash(input string) string {
	args := m.Called(input)
	return args.String(0)
}

func TestShortenerService(t *testing.T) {
	ctx := context.Background()

	t.Run("Create data successfully", func(t *testing.T) {
		mockRepo := new(MockShortURLRepository)
		mockHasher := new(MockHasher)

		originalURL := "https://example.com/long/url/1"
		expectedHash := "hash_of_long_url_1"

		mockHasher.On("Hash", originalURL).Return(expectedHash).Once()

		mockRepo.On("FindByURLHash", ctx, expectedHash).Return([]*domain.ShortURL{}, nil).Once()

		mockRepo.On("Insert", ctx, mock.AnythingOfType("*domain.ShortURL")).Return(uint64(1), nil).Once()

		mockRepo.On("UpdateShortCode", ctx, uint64(1), "1").Return(nil).Once()

		service := app.NewShortenerService(mockRepo, mockHasher)

		shortCode, err := service.ShortenURL(ctx, originalURL)

		assert.NoError(t, err)
		assert.Equal(t, "1", shortCode)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("Create existing data", func(t *testing.T) {
		mockRepo := new(MockShortURLRepository)
		mockHasher := new(MockHasher)

		originalURL := "https://example.com/long/url/2"
		existingHash := "hash_of_long_url_2"
		existingShortCode := "xyz789"

		existingShortURL := &domain.ShortURL{
			ID:          2,
			OriginalURL: originalURL,
			URLHash:     existingHash,
			ShortCode:   existingShortCode,
		}

		mockHasher.On("Hash", originalURL).Return(existingHash).Once()

		mockRepo.On("FindByURLHash", ctx, existingHash).Return([]*domain.ShortURL{existingShortURL}, nil).Once()

		mockRepo.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything)
		mockRepo.AssertNotCalled(t, "UpdateShortCode", mock.Anything, mock.Anything, mock.Anything)

		service := app.NewShortenerService(mockRepo, mockHasher)

		shortCode, err := service.ShortenURL(ctx, originalURL)

		assert.NoError(t, err)
		assert.Equal(t, existingShortCode, shortCode)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("Create data with hash collision but different original URL", func(t *testing.T) {
		mockRepo := new(MockShortURLRepository)
		mockHasher := new(MockHasher)

		originalURL1 := "https://example.com/long/url/collision_a"
		originalURL2 := "https://example.com/long/url/collision_b"
		collidingHash := "colliding_hash_value"

		existingCollidingURL := &domain.ShortURL{
			ID:          3,
			OriginalURL: originalURL1,
			URLHash:     collidingHash,
			ShortCode:   "xyzabc",
		}

		mockHasher.On("Hash", originalURL2).Return(collidingHash).Once()
		mockRepo.On("FindByURLHash", ctx, collidingHash).Return([]*domain.ShortURL{existingCollidingURL}, nil).Once()
		mockRepo.On("Insert", ctx, mock.AnythingOfType("*domain.ShortURL")).Return(uint64(4), nil).Once()
		mockRepo.On("UpdateShortCode", ctx, uint64(4), "4").Return(nil).Once()

		service := app.NewShortenerService(mockRepo, mockHasher)

		shortCode, err := service.ShortenURL(ctx, originalURL2)

		assert.NoError(t, err)
		assert.Equal(t, "4", shortCode)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})
}
