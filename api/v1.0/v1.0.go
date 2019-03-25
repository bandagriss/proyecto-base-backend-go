package apiv1

import (
	"github.com/gin-gonic/gin"
	"./auth"
	"./posts"
)

// ping Prueba ping pong
func ping(c *gin.Context)  {
  c.JSON(200, gin.H{
		"message": "pong",
	})
}

// ApplyRoutes ...
func ApplyRoutes(r *gin.RouterGroup)  {
  v1 := r.Group("/v1.0")
	{
		v1.GET("/ping", ping)
		auth.ApplyRoutes(v1)
		posts.ApplyRoutes(v1)
	}
}
