package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, postHandler *PostHandler) {
	api := r.Group("/api/v1")
	{
		api.GET("/posts", postHandler.GetAll)
		api.GET("/posts/:slug", postHandler.Get)
		api.POST("/posts/:slug/view", postHandler.IncreaseView)
		api.POST("/posts/:slug/like", postHandler.IncreaseLike)
	}
}
