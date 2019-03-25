package api

import (
	"github.com/gin-gonic/gin"
	"./v1.0"
)

// ApplyRoutes aplicando router a gin Router
func ApplyRoutes(r *gin.Engine)  {
  api := r.Group("/api")
	{
		apiv1.ApplyRoutes(api)
	}
}
