package middleware

import (
	"github.com/gin-gonic/gin"
)

func ConfigOperatorName(operatorName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("operatorName", operatorName)
		c.Next()
	}
}
