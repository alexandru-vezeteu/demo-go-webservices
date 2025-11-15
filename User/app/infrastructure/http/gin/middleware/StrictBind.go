package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
)

func StrictBindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		if err == io.EOF {
			return fmt.Errorf("request body cannot be empty")
		}
		return fmt.Errorf("invalid request format: %v", err)
	}
	return nil
}
