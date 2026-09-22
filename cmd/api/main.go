package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shikihtm/blog-backend/internal/database"
	"github.com/shikihtm/blog-backend/internal/handler/auth"
	handler "github.com/shikihtm/blog-backend/internal/handler/post"
	"github.com/shikihtm/blog-backend/internal/repository"
)

func main() {
	log.Println("[MAIN] [INFO] Starting Shiki Blog Backend service...")

	dbConn, err := database.Initialize("/app/data/blog.db")
	if err != nil {
		log.Fatalf("[MAIN] [FATAL] Database initialization failed: %v", err)
	}
	defer dbConn.Close()

	authConfig, err := auth.LoadConfig()
	if err != nil {
		log.Fatalf("[MAIN] [FATAL] Config load failed: %v", err)
	}

	repo := repository.NewRepository(dbConn)
	postHandler := handler.NewPostHanlder(repo, authConfig.JWTSecretKey)
	authHandler := auth.NewAuthenticationHandler(authConfig)

	repository.SyncAll(repo)
	repository.Watch(repo)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://shikii.dev"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r := router.Group("/api/v1")

	postHandler.RegisterRoutes(r)
	authHandler.RegisterRoutes(r)

	log.Println("[MAIN] [INFO] HTTP Server is ready on port :3000")
	if err := router.Run(":3050"); err != nil {
		log.Fatalf("[MAIN] [FATAL] Failed to start HTTP server: %v", err)
	}
}
