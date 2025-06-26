package web

import "github.com/gin-gonic/gin"

func SetupRouter(handler *ShortenerHandler) *gin.Engine {
	r := gin.Default()

	r.POST("/api/shortlinks", handler.ShortenURLHandler)
	r.GET("/api/shortlinks/:shortID", handler.RedirectURLHandler)

	return r
}
