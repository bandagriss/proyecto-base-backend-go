package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Authorized blocks unauthorized requests
func Authorized(c *gin.Context) {
	_, exists := c.Get("user")
	if !exists {
		fmt.Println("Ocurrio un error al autorizar ===> aqui ", exists)
		c.AbortWithStatus(401)
		return
	}
}
