package main

import (
	"beaver/inventory/adapters/eventstore"
	"beaver/inventory/adapters/http"
	"beaver/inventory/auth"
	"beaver/inventory/config"
	"beaver/inventory/core/domain"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		panic("no port defined")
	}

	dbConn := os.Getenv("DB")

	if dbConn == "" {
		panic("no db connection string")
	}

	IdpUrl := os.Getenv("IDP_URL")

	if IdpUrl == "" {
		panic("no IdpUrl defined")
	}

	cfg := config.LoadConfig(dbConn)
	config.InitDb(*cfg)
	defer cfg.DB.Close()
	eventStore := eventstore.NewPostgresEventStore(cfg.DB)

	stockService := domain.NewStockService(eventStore)

	err := stockService.RebuildEventStream() // Event Stream beim Start neu bilden
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	r.Use(cors.New(config))
	r.Use(auth.TokenCheck(IdpUrl))
	v1 := r.Group("/v1")

	http.NewV1Handler(v1, stockService)

	r.Run(":" + port)

}
