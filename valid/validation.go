package valid

import "fmt"
import "github.com/picatic/go-api/norm"


type Validator interface {
	Validate(interface{}) (bool, error)
}


// Validatable
type Validatable interface {
	LoadValidations() // loads validation rules into Fields
	ValidationsLoaded() bool // have we loaded the valudation rules into Fields
	Validators() Validators
}

type ValidatableModel interface {
	Validatable
	norm.Model
}

type Validators map[string]map[string]ValidatorFunc

// Simple
// Complex with args
// Complex Async with args
type ValidatorFunc func(model norm.Model, field string, args... interface{}) (<-chan error)

// validate all fields
func Validate(model ValidatableModel) (<-chan error) {
	err := make(chan error, 1)
	err <- nil
	return err
}

// validates indicated fields
func ValidateFields(model ValidatableModel, fields[]string) (<-chan error) {
	err := make(chan error, 1)
	err <- nil
	return err
}

func AddFieldValidation(validatable Validatable, fieldName string, alias string, validator Validator) {

}
func RemoveFieldValidation(validatable Validatable, fieldName string, alias string) {

}

type ValidationError struct {
	Field string
	Message string
	Model norm.Model
}

func (ve *ValidationError) Error() string{
	return fmt.Sprintf("Field: %s Error: %s", ve.Field, ve.Message)
}

type ValidationModelError struct {
	Message string
	Model norm.Model
}

func (vme *ValidationModelError) Error() string {
	return fmt.Sprintf("Model %+v", vme.Model)
}

type ValidationErrors struct {
	Errors []error
	Model norm.Model
}

func (ve *ValidationErrors) Error() string {
	return "NOPE"
}