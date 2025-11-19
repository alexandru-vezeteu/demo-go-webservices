package middleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"userService/application/domain"
)

func ParseIDParam(c *gin.Context, paramName string) (int, error) {
	idStr := c.Param(paramName)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, &domain.ValidationError{Field: paramName, Reason: "must be an integer"}
	}
	if id <= 0 {
		return 0, &domain.ValidationError{Field: paramName, Reason: "must be greater than 0"}
	}
	return id, nil
}
