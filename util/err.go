package util

import (
	"github.com/go-playground/validator/v10"
)

func ListOfErrors(t interface{}, e error) []map[string]string {
	ve := e.(validator.ValidationErrors)
	InvalidFields := make([]map[string]string, 0)

	for _, e := range ve {
		errors := map[string]string{}
		errors[e.Field()] = msgForTag(e.Tag())
		InvalidFields = append(InvalidFields, errors)
	}

	return InvalidFields
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "min":
		return "This field is below the minimum required length"
	case "max":
		return "This field is above the maximum allowed length"
	}
	return ""
}
