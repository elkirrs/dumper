package validation

import "github.com/go-playground/validator/v10"

func xorValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	otherFieldName := fl.Param()
	other := fl.Parent().FieldByName(otherFieldName)

	if !other.IsValid() {
		return false
	}

	a := field.String()
	b := other.String()

	return a != "" || b != ""
}
