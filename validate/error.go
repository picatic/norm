package validate

import (
	"errors"
)

var (
	ErrNotString = errors.New("value is not a string")
)

type ValidationError struct {
	Alias string
	Locationer
	Err string
}

func (ve ValidationError) Error() string {
	return ve.Location() + ": " + ve.Err
}

func (ve *ValidationError) AddLocation(l Locationer) {
	switch location := ve.Locationer.(type) {
	case nil:
		ve.Locationer = l
	case EmbeddedLocation:
		ve.Locationer = EmbeddedLocation(append([]Locationer{l}, location...))
	default:
		ve.Locationer = EmbeddedLocation(append([]Locationer{l}, location))
	}
}

type ValidationErrors []*ValidationError

func (ves ValidationErrors) Error() string {
	errorString := ves[0].Error()

	for i := 1; i < len(ves); i++ {
		errorString += ", " + ves[i].Error()
	}

	return errorString
}

func (ves *ValidationErrors) AddError(err error) {
	switch err := err.(type) {
	case ValidationError:
		*ves = append(*ves, &err)
	case ValidationErrors:
		*ves = append(*ves, err...)
	default:
		e := NewError(err.Error())
		*ves = append(*ves, &e)
	}
}

func (ves *ValidationErrors) AddLocation(l Locationer) {
	for _, ve := range *ves {
		ve.AddLocation(l)
	}
}

func Alias(alias string, validator Validator) Validator {
	return ValidatorFunc(func(v interface{}) error {
		err := validator.Validate(v)
		if err != nil {
			ve := err.(ValidationError) //should only Alias single error otherwise panic
			if ve.Alias != "" {
				panic("error should not be aliased already")
			}
			ve.Alias = alias
			return ve
		}

		return nil
	})
}

func NewError(err string) ValidationError {
	return ValidationError{
		Err: err,
	}
}
