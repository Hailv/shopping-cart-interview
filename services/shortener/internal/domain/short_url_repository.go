package domain

import "context"

type ShortURLRepository interface {
	Insert(ctx context.Context, shortURL *ShortURL) (uint64, error)
	UpdateShortCode(ctx context.Context, id uint64, shortCode string) error
	FindByShortCode(ctx context.Context, shortCode string) (*ShortURL, error)
	FindByURLHash(ctx context.Context, urlHash string) ([]*ShortURL, error)
	FindAll(ctx context.Context) ([]*ShortURL, error)
}
