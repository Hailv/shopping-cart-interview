package web

import (
	"context"
	"errors"
	"net/http"

	"github.com/cinchprotocol/cinch-api/services/shortener/internal/app"
	"github.com/gin-gonic/gin"
)

type ShortenerServicePort interface {
	ShortenURL(ctx context.Context, originalURL string) (string, error)   // Trả về short code
	GetOriginalURL(ctx context.Context, shortCode string) (string, error) // Nhận short code
}

type ShortenRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"` // Sẽ chứa short_code
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type ShortenerHandler struct {
	service ShortenerServicePort
}

func NewShortenerHandler(svc ShortenerServicePort) *ShortenerHandler {
	return &ShortenerHandler{service: svc}
}

func (h *ShortenerHandler) ShortenURLHandler(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: app.ErrInvalidInput.Error()})
		return
	}

	shortCode, err := h.service.ShortenURL(c.Request.Context(), req.OriginalURL)
	if err != nil {
		if errors.Is(err, app.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error when get url"})
		return
	}

	c.JSON(http.StatusCreated, ShortenResponse{ShortURL: shortCode})
}

func (h *ShortenerHandler) RedirectURLHandler(c *gin.Context) {
	shortCode := c.Param("shortID")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: app.ErrInvalidInput.Error()})
		return
	}

	originalURL, err := h.service.GetOriginalURL(c.Request.Context(), shortCode)
	if err != nil {
		if errors.Is(err, app.ErrURLNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error when get url"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}
