package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateRequest(schema interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(schema); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("validatedData", schema)
		c.Next()
	}
}
