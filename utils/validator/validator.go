package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	ValidationErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}
	XValidator struct {
		validator *validator.Validate
	}
)

var Validator = &XValidator{validator: validator.New()}

func (v XValidator) Validate(data interface{}) []ValidationErrorResponse {
	validationErrors := []ValidationErrorResponse{}

	errs := v.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ValidationErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

func ParseAndValidate(c *fiber.Ctx, target interface{}) []ValidationErrorResponse {
	if err := c.BodyParser(target); err != nil {
		return []ValidationErrorResponse{
			{
				Error:       true,
				FailedField: "body",
				Tag:         err.Error(),
			},
		}
	}
	if errArr := Validator.Validate(target); len(errArr) != 0 {
		return errArr
	}
	return nil
}
