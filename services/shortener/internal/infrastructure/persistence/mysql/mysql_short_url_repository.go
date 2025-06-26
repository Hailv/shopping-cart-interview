package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/cinchprotocol/cinch-api/services/shortener/internal/app"
	"github.com/cinchprotocol/cinch-api/services/shortener/internal/domain"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type MySQLShortURLRepository struct {
	db *gorm.DB
}

func NewMySQLShortURLRepository(dataSourceName string) (*MySQLShortURLRepository, error) {
	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to mysql via GORM: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to mysql via GORM: %w", err)
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return &MySQLShortURLRepository{db: db}, nil
}

func (r *MySQLShortURLRepository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
func (r *MySQLShortURLRepository) Insert(ctx context.Context, shortURL *domain.ShortURL) (uint64, error) {
	result := r.db.WithContext(ctx).Create(shortURL)
	if result.Error != nil {
		return 0, result.Error
	}
	return shortURL.ID, nil
}

func (r *MySQLShortURLRepository) UpdateShortCode(ctx context.Context, id uint64, shortCode string) error {
	result := r.db.WithContext(ctx).Model(&domain.ShortURL{}).Where("id = ?", id).Update("short_code", shortCode)
	return result.Error
}

func (r *MySQLShortURLRepository) FindByShortCode(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	var shortURL domain.ShortURL
	result := r.db.WithContext(ctx).Where("short_code = ?", shortCode).First(&shortURL)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app.ErrURLNotFound
		}
		return nil, result.Error
	}
	return &shortURL, nil
}

func (r *MySQLShortURLRepository) FindByURLHash(ctx context.Context, urlHash string) ([]*domain.ShortURL, error) {
	var foundURLs []*domain.ShortURL
	result := r.db.WithContext(ctx).Where("url_hash = ?", urlHash).Find(&foundURLs)
	fmt.Println(result)
	if result.Error != nil {
		return nil, result.Error
	}
	return foundURLs, nil
}

func (r *MySQLShortURLRepository) FindAll(ctx context.Context) ([]*domain.ShortURL, error) {
	var allURLs []*domain.ShortURL
	result := r.db.WithContext(ctx).Find(&allURLs)
	if result.Error != nil {
		return nil, result.Error
	}
	return allURLs, nil
}
