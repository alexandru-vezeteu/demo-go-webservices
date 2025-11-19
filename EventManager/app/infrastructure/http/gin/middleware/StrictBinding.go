package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"eventManager/application/domain"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// StrictBindJSON binds JSON with strict validation:
// - Rejects unknown fields
// - Rejects duplicate keys
// - Returns detailed error messages
func StrictBindJSON(c *gin.Context, obj interface{}) error {
	// Read the body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return &domain.InvalidRequestError{Reason: "failed to read request body"}
	}

	// Restore the body for potential future reads
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Check for empty body
	if len(body) == 0 {
		return &domain.InvalidRequestError{Reason: "request body cannot be empty"}
	}

	// Create a decoder with DisallowUnknownFields
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()

	// Decode with strict validation
	if err := decoder.Decode(obj); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return &domain.InvalidRequestError{Reason: fmt.Sprintf("malformed JSON at position %d", syntaxError.Offset)}
		case errors.As(err, &unmarshalTypeError):
			return &domain.InvalidRequestError{Reason: fmt.Sprintf("invalid value for field '%s' (expected %s)", unmarshalTypeError.Field, unmarshalTypeError.Type)}
		case errors.Is(err, io.EOF):
			return &domain.InvalidRequestError{Reason: "request body is empty"}
		case errors.Is(err, io.ErrUnexpectedEOF):
			return &domain.InvalidRequestError{Reason: "malformed JSON"}
		default:
			// Check for unknown field error (safely)
			errMsg := err.Error()
			if strings.HasPrefix(errMsg, "json: unknown field") {
				return &domain.InvalidRequestError{Reason: errMsg}
			}
			return &domain.InvalidRequestError{Reason: err.Error()}
		}
	}

	// Check if there's additional data after the JSON object
	if decoder.More() {
		return &domain.InvalidRequestError{Reason: "request body contains multiple JSON objects"}
	}

	return nil
}

// ParseIDParam validates and parses an integer ID from route parameters
func ParseIDParam(c *gin.Context, paramName string) (int, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		return 0, &domain.ValidationError{Field: paramName, Reason: "parameter is required"}
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, &domain.ValidationError{Field: paramName, Reason: "must be an integer"}
	}

	if id < 1 {
		return 0, &domain.ValidationError{Field: paramName, Reason: "must be greater than 0"}
	}

	return int(id), nil
}

// StrictBindQuery validates query parameters and rejects unknown parameters
func StrictBindQuery(c *gin.Context, obj interface{}, allowedParams []string) error {
	// First check for unknown query parameters
	queryParams := c.Request.URL.Query()
	allowedMap := make(map[string]bool)
	for _, param := range allowedParams {
		allowedMap[param] = true
	}

	for param := range queryParams {
		if !allowedMap[param] {
			return &domain.InvalidRequestError{Reason: fmt.Sprintf("unknown query parameter: %s", param)}
		}
	}

	// Then bind the query parameters
	if err := c.ShouldBindQuery(obj); err != nil {
		return &domain.InvalidRequestError{Reason: err.Error()}
	}

	return nil
}

// RejectUnknownJSONMiddleware is a middleware that validates JSON requests globally
func RejectUnknownJSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodDelete {
			contentType := c.ContentType()
			if contentType != "" && contentType != "application/json" {
				c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
					"error": "Content-Type must be application/json",
				})
				return
			}
		}
		c.Next()
	}
}
