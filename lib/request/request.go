package request

import (
	"errors"
	"github.com/Verce11o/Hezzl-Go/lib/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func CustomTagNameFunc() validator.TagNameFunc {
	return func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	}
}

func ValidationErrors(errs validator.ValidationErrors) gin.H {
	errMessages := gin.H{}
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessages[err.Field()] = "required"
		case "min":
			errMessages[err.Field()] = "shorter than " + err.Param()
		case "max":
			errMessages[err.Field()] = "longer than " + err.Param()
		default:
			errMessages[err.Field()] = "not valid"
		}
	}

	return errMessages
}

func Read(c *gin.Context, request any) any {
	if err := c.ShouldBindJSON(request); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			return ValidationErrors(ve)
		}
		return response.ErrInvalidRequest.Error()
	}
	return nil

}
