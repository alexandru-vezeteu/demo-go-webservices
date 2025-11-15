package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ParseIDParam(c *gin.Context, paramName string) (int, error) {
	idStr := c.Param(paramName)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: must be an integer", paramName)
	}
	if id <= 0 {
		return 0, fmt.Errorf("invalid %s: must be greater than 0", paramName)
	}
	return id, nil
}
