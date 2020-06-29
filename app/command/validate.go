package command

import (
	"gopkg.in/go-playground/validator.v9"
)

// use a single instance of Validate, it caches struct info
var (
	validate = validator.New()
)

func validateOptions(value interface{}) error {
	return validate.Struct(value)
}
