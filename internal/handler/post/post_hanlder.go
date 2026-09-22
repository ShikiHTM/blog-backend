package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/shikihtm/blog-backend/internal/model"
	"github.com/shikihtm/blog-backend/internal/repository"
	"github.com/shikihtm/blog-backend/middleware"
)

type PostHandler struct {
	repo      repository.BlogRepository
	jwtSecret string
}

func NewPostHanlder(r repository.BlogRepository, jwtSecret string) *PostHandler {
	return &PostHandler{
		repo:      r,
		jwtSecret: jwtSecret,
	}
}

func (h *PostHandler) GetAll(c *gin.Context) {
	posts, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"error":   "INTERNAL_SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) Get(c *gin.Context) {
	slug := c.Param("slug")

	post, err := h.repo.Get(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Post not found",
			"error":   "ERR_NOT_FOUND",
		})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req model.PostCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   "ERR_BAD_REQUEST",
		})
		return
	}

	safeSlug := filepath.Base(filepath.Clean(req.Slug))
	if safeSlug == "." || safeSlug == "/" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid slug",
			"error":   "ERR_BAD_REQUEST",
		})
		return
	}

	fileContent := fmt.Sprintf(`---
title: %s
topic: %s
cover: %s
author: %s
---

%s
	`, req.Title, req.Topic, req.Cover, req.Author, req.Content)

	savePath := filepath.Join(".", "posts", fmt.Sprintf("%s.mdx", safeSlug))

	if err := os.WriteFile(savePath, []byte(fileContent), 0644); err != nil {
		fmt.Printf("Error: Cannot write to file: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"error":   "INTERNAL_SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Post created successfully!",
	})
}

func (h *PostHandler) IncreaseLike(c *gin.Context) {
	slug := c.Param("slug")

	_, err := h.repo.IncreaseLike(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"error":   "INTERNAL_SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *PostHandler) IncreaseView(c *gin.Context) {
	slug := c.Param("slug")

	_, err := h.repo.IncreaseView(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"error":   "INTERNAL_SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *PostHandler) RegisterRoutes(r *gin.RouterGroup) {
	posts := r.Group("/posts")
	{
		posts.GET("/", h.GetAll)
		posts.GET("/:slug", h.Get)
		posts.POST("/", h.CreatePost).Use(middleware.RequireAuthCookie(h.jwtSecret))
	}
}
