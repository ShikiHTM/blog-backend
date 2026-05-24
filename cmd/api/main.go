package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shikihtm/blog-backend/internal/database"
	"github.com/shikihtm/blog-backend/internal/handler"
	"github.com/shikihtm/blog-backend/internal/repository"
)

func main() {
	log.Println("[MAIN] [INFO] Starting Shiki Blog Backend service...")

	dbConn, err := database.Initialize("./data/blog.db")
	if err != nil {
		log.Fatalf("[MAIN] [FATAL] Database initialization failed: %v", err)
	}
	defer dbConn.Close()

	repo := repository.NewRepository(dbConn)
	postHandler := handler.NewPostHanlder(repo)

	repository.SyncAll(repo)
	repository.Watch(repo)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	handler.RegisterRoutes(router, postHandler)

	log.Println("[MAIN] [INFO] HTTP Server is ready on port :3000")
	if err := router.Run(":3050"); err != nil {
		log.Fatalf("[MAIN] [FATAL] Failed to start HTTP server: %v", err)
	}
}
