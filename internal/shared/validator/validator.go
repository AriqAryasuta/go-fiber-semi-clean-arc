package validator

import "github.com/go-playground/validator/v10"

type Validator struct {
	engine *validator.Validate
}

func New() *Validator {
	return &Validator{engine: validator.New()}
}

func (v *Validator) ValidateStruct(payload any) error {
	return v.engine.Struct(payload)
}
