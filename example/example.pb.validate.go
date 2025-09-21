package example

import (
	validator "github.com/go-playground/validator/v10"
)

func (m *Example) Validate() error {
	return validator.New().Struct(m)
}
