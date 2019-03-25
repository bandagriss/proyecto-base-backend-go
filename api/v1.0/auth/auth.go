package auth

import (
	"github.com/gin-gonic/gin"
)


// ApplyRoutes aplicando router a gin Engine
func ApplyRoutes(r *gin.RouterGroup)  {
  auth := r.Group("/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
		auth.GET("/check", check)
	}
}
