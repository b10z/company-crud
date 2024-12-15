package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	*validator.Validate
}

// Some more opts/configs could be added here.
func New() *Validator {
	vldtr := validator.New()

	return &Validator{vldtr}
}
