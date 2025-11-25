package validation

import (
	"dumper/internal/domain/config"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validation struct {
	validator *validator.Validate
}

func New() *Validation {
	v := validator.New()

	v.RegisterValidation("xor", xorValidator)
	v.RegisterTagNameFunc(fieldNameExtractor)

	return &Validation{validator: v}
}

func HumanError(err error) error {
	if err == nil {
		return nil
	}

	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		var messages []string

		for _, e := range errs {
			switch e.Tag() {
			case "xor":
				a := e.Field()
				b := e.Param()
				messages = append(messages, fmt.Sprintf("either %s or %s must be set", a, b))
			default:
				messages = append(messages, fmt.Sprintf("%s is invalid (%s)", e.Field(), e.Tag()))
			}
		}

		if len(messages) > 0 {
			return fmt.Errorf("%s", joinMessages(messages))
		}
	}

	return err
}

func joinMessages(msgs []string) string {
	return fmt.Sprintf("%s", msgs)
}

func (v *Validation) Handler(cfg *config.Config) error {
	if err := v.validator.Struct(cfg); err != nil {
		return err
	}

	if err := validateServer(v, cfg); err != nil {
		return err
	}

	if err := validateDatabase(v, cfg); err != nil {
		return err
	}

	if err := validateStorages(v, cfg); err != nil {
		return err
	}

	return nil
}
