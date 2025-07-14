package utils

import (
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()
}

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	for _, fieldErr := range err.(validator.ValidationErrors) {
		errors[fieldErr.Field()] = validationMessage(fieldErr)
	}
	return errors
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Обязательное поле"
	case "email":
		return "Неверный формат email"
	case "min":
		return "Минимум " + fe.Param() + " символов"
	default:
		return "Неверное значение"
	}
}
