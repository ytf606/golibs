package validate

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type (
	ValidationErrors = validator.ValidationErrors
	FieldError       = validator.FieldError
)

type ValidError struct {
	Key string
	Msg string
}

type ValidErrors []*ValidError

func (v *ValidError) Error() string {
	return v.Msg
}

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

func ToError(err interface{}) ValidErrors {
	var errs ValidErrors

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return errs
	}

	for key, value := range verrs.Translate(trans) {
		errs = append(errs, &ValidError{
			Key: key,
			Msg: value,
		})
	}

	return errs
}
