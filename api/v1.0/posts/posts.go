package posts

import (
	"github.com/gin-gonic/gin"
	"../../../lib/middlewares"
)

func ApplyRoutes(r *gin.RouterGroup) {
	posts := r.Group("/posts")
	{
		posts.POST("/", middlewares.Authorized, create)
		posts.GET("/", list)
		posts.GET("/:id", read)
		posts.DELETE("/:id", middlewares.Authorized, remove)
		posts.PATCH("/:id", middlewares.Authorized, update)
	}
}

