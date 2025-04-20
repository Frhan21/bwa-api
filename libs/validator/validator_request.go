package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	var errMessage []string
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "email":
				errMessage = append(errMessage, "Invalid email format")
			case "required":
				errMessage = append(errMessage, "Field "+err.Field()+" is required")
			case "min":
				if err.Field() == "Password" {
					errMessage = append(errMessage, "Password must have at least 8 characters")
				}
			case "eqfield":
				errMessage = append(errMessage, "Field "+err.Field()+" must be equal to "+err.Param())
			default:
				errMessage = append(errMessage, "Field"+err.Field()+"is not valid")
			}
		}

		return errors.New("Validasi Gagal: " + joinMessage(errMessage))
	}

	return nil
}

func joinMessage(messages []string) string {
	result := ""
	for i, message := range messages {
		if i > 0 {
			result += ", "
		}
		result += message
	}

	return result
}
