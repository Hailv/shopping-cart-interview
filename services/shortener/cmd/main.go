package main

import (
	"fmt"
	"github.com/cinchprotocol/cinch-api/services/shortener/internal/app"
	"github.com/cinchprotocol/cinch-api/services/shortener/internal/infrastructure/persistence/mysql"
	"github.com/cinchprotocol/cinch-api/services/shortener/internal/infrastructure/web"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("get file .env has error %v", err)
	}

	dbUser := "root"
	dbPassword := os.Getenv("DB_ROOT_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbParams := os.Getenv("DB_PARAMS")

	dbConnStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	if dbParams != "" {
		dbConnStr = fmt.Sprintf("%s?%s", dbConnStr, dbParams)
	}

	repo, err := mysql.NewMySQLShortURLRepository(dbConnStr)
	if err != nil {
		log.Fatalf("Cannot connect MySQL: %v", err)
	}
	defer func() {
		if err := repo.Close(); err != nil {
			log.Printf("Cannot connect DB: %v", err)
		}
	}()

	shortenerService := app.NewShortenerService(repo, &app.SHA256Hasher{})
	handler := web.NewShortenerHandler(shortenerService)
	router := gin.Default()
	router = web.SetupRouter(handler)

	// Chạy server.
	log.Println("server listen port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
