package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"userService/application/domain"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func StrictBindJSON(c *gin.Context, obj interface{}) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return &domain.InvalidRequestError{Reason: "failed to read request body"}
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if len(body) == 0 {
		return &domain.InvalidRequestError{Reason: "request body cannot be empty"}
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()

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
			errMsg := err.Error()
			if strings.HasPrefix(errMsg, "json: unknown field") {
				return &domain.InvalidRequestError{Reason: errMsg}
			}
			return &domain.InvalidRequestError{Reason: err.Error()}
		}
	}

	if decoder.More() {
		return &domain.InvalidRequestError{Reason: "request body contains multiple JSON objects"}
	}


	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.Struct(obj); err != nil {
			var validationErrs validator.ValidationErrors
			if errors.As(err, &validationErrs) {

				for _, fieldErr := range validationErrs {
					return &domain.ValidationError{
						Field:  fieldErr.Field(),
						Reason: fmt.Sprintf("validation failed on '%s' tag", fieldErr.Tag()),
					}
				}
			}
			return &domain.ValidationError{Reason: err.Error()}
		}
	}

	return nil
}

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

func StrictBindQuery(c *gin.Context, obj interface{}, allowedParams []string) error {
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

	if err := c.ShouldBindQuery(obj); err != nil {
		return &domain.InvalidRequestError{Reason: err.Error()}
	}

	return nil
}
