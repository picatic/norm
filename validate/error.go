package validate

import (
	"errors"
)

var (
	ErrNotString = errors.New("value is not a string")
)

type ValidationError struct {
	Alias *string
	Locationer
	Err error
}

func (ve ValidationError) Error() string {
	return ve.Location() + ": " + ve.Err.Error()
}

type ValidationErrors []ValidationError

func (ves ValidationErrors) Error() string {
	errorString := ves[0].Error()

	for i := 1; i < len(ves); i++ {
		errorString += ", " + ves[i].Error()
	}

	return errorString
}
