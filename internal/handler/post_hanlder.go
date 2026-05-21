package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shikihtm/blog-backend/internal/repository"
)

type PostHandler struct {
	repo repository.BlogRepository
}

func NewPostHanlder(r repository.BlogRepository) *PostHandler {
	return &PostHandler{
		repo: r,
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
