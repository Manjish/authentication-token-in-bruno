package utils

import (
	"bruno_authentication/pkg/api_errors"
	"bruno_authentication/pkg/framework"
	"errors"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func HandleValidationError(logger framework.Logger, ctx *gin.Context, err error) {
	logger.Error(err)
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}

type ErrorMsg struct {
	Field     string `json:"field"`
	FullField string `json:"full_field"`
	Message   string `json:"message"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "max":
		return "Should be less than " + fe.Param()
	case "min":
		return "Should be greater than " + fe.Param()
	}
	return "Unknown error"
}

func HandleValidationWithError(logger framework.Logger, c *gin.Context, err error) {
	logger.Error("VALIDATION ERROR:", err)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ErrorMsg, len(ve))
		for i, fe := range ve {
			out[i] = ErrorMsg{
				fe.Field(),
				fe.Namespace(),
				getErrorMsg(fe),
			}
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
}

// list static errors to filter
var exceptStaticError = []error{
	gorm.ErrRecordNotFound,
	api_errors.ErrInvalidUUID,
}

// list dyanmic errors to filter
var exceptDynamicError = []error{}

// list SQL errors to filter
var exceptSQLError = []uint16{
	1062, // duplicate entry
}

var sqlError *mysql.MySQLError

func HandleError(logger framework.Logger, c *gin.Context, err error) {
	logger.Error(err)

	// will not captured by sentry if its an explicit APIError
	if apiErr, ok := err.(*api_errors.APIError); ok {
		c.JSON(apiErr.StatusCode, gin.H{
			"error": apiErr.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})

	for _, e := range exceptStaticError {
		if errors.Is(err, e) {
			return
		}
	}

	for _, e := range exceptDynamicError {
		if reflect.TypeOf(e) == reflect.TypeOf(err) {
			return
		}
	}

	if errors.As(err, &sqlError) {
		for _, code := range exceptSQLError {
			if code == sqlError.Number {
				return
			}
		}
	}

	sentry.CaptureException(err)
}
