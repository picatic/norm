package norm

import (
	"errors"
	"fmt"
	"github.com/picatic/norm/field"
	"reflect"
)

// Validators global static map of validators for models
var Validators = ValidatorMap{}

// Validator Models can implement a custom validator, which is expected to return ValidationError/ValidationErrors
type Validator interface {
	Validate(interface{}) error
}

// ValidatorMap Store validators for models
type ValidatorMap map[reflect.Type][]Validator

// Get Fetch the validators registered to a Model
func (vm ValidatorMap) Get(model interface{}) []Validator {
	modelType := reflect.TypeOf(model)
	validators, ok := vm[modelType]
	if !ok {
		validators = []Validator{}
		vm[modelType] = validators
	}
	return validators
}

// Set a validator
func (vm ValidatorMap) Set(modelType reflect.Type, validators []Validator) {
	vm[modelType] = validators
}

// Del Delete a validator
func (vm ValidatorMap) Del(modelType reflect.Type) {
	delete(vm, modelType)
}

// Clone the validators
func (vm ValidatorMap) Clone() ValidatorMap {
	newMap := ValidatorMap{}
	for k, v := range vm {
		newMap[k] = []Validator{}
		copy(newMap[k], v)
	}
	return newMap
}

// Validate a model and specified fields.
func (vm ValidatorMap) Validate(model interface{}, fields field.Names) *ValidationErrors {
	errs := &ValidationErrors{}
	for _, validator := range vm.Get(model) {
		switch v := validator.(type) {
		case FieldValidator:
			for _, field := range fields {
				if field == v.Field {
					if err := v.Validate(model); err != nil {
						errs.Add(err)
						break
					}
				}
			}
		}
	}
	if len(errs.Errors) > 0 {
		return errs
	}
	return nil
}

// FieldValidator Wrapper for a standard field validation
type FieldValidator struct {
	Field field.Name
	Alias string
	Func  FieldValidatorFunc
	Args  []interface{}
}

// Validate a field on a model
func (fv FieldValidator) Validate(model interface{}) error {
	return fv.Func(model, fv.Field, fv.Args)
}

// NewFieldValidator Create a new FieldValidator, providing a FieldName, alias for the error,
// a function to run to validate and any args required for that validator.
func NewFieldValidator(
	field field.Name,
	alias string,
	vFunc FieldValidatorFunc,
	args ...interface{},
) Validator {
	return &FieldValidator{
		Field: field,
		Alias: alias,
		Func:  vFunc,
		Args:  args,
	}
}

// NewGovalidator FieldValidator based on go validators
func NewGovalidator(
	field field.Name,
	alias string,
	vFunc StringValidator,
) Validator {
	return &FieldValidator{
		Field: field,
		Alias: alias,
		Func:  NewStringValidatorFunc(vFunc),
		Args:  make([]interface{}, 0),
	}
}

//func (v Validator) String() string {
//	return fmt.Sprintf("%s %s %s(%s)", v.Field, v.Alias, "func", v.Args)
//}

//func NewValidator(field field.Name, alias string, fn ValidatorFunc, args ...interface{}) Validator {
//	return Validator{field, alias, fn, args}
//}

//func (fv *FieldValidator) Validate(m Model, fields []field.Name) error {
//	return fv.Func(m, field, fv.Args...)
//}
//
//
//type ModelValidator struct {
//	Func ModelValidatorFunc
//	Args []interface{}
//}
//
//func (mv *ModelValidator) Validate(m Model, fields []field.Name) error {
//	return mv.Func(m, fields, mv.Args...)
//}
//

// Validatable
//type Validatable interface {
//	ValidatorList() ValidatorList // Must be locking
//	ValidatorListReset()          // Reset ValidatorMap
//}

//type ValidatableModel interface {
//	Validatable
//	Model
//}

//type ValidatorList []Validator
//
//// Add a validator
//func (vl *ValidatorList) AddValidator(v Validator) {
//	vl = append(vl, v)
//}
//
// Validate a model with this map for all fields indicated. Nil returned if no validator errors
//func (vl *ValidatorList) Validate(model Model, fields field.Names) <-chan error {
//	err := make(chan error, 1)
//	err <- nil
//	return err
//}

// Simple
// Complex with args
// Complex Async with args

// FieldValidatorFunc What a validator must implement
type FieldValidatorFunc func(value interface{}, args ...interface{}) error

//var vMap ValidatorMap = &ValidatorMap{"Id":{"IsString": IsString}}

//ValidationError Represent a single validation error
type ValidationError struct {
	Field   string
	Message string
	Model   Model
}

// Error String the error
func (ve ValidationError) Error() string {
	return fmt.Sprintf("Field: %s Error: %s", ve.Field, ve.Message)
}

// ValidationErrors Represent a set of ValidationError
type ValidationErrors struct {
	Errors []error
	Model  Model
}

// Add an error to this set
func (ve *ValidationErrors) Add(err error) {
	ve.Errors = append(ve.Errors, err)
}

// Error string the error
func (ve ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return fmt.Sprintf("Empty errors")
	}
	if fe, ok := ve.Errors[0].(*ValidationError); ok == true {
		return fmt.Sprintf("First of multiple errors, Field: %s Error: %s", fe.Field, fe.Message)
	}
	return fmt.Sprintf("First of multiple errors is not a Validation error")
}

// AddValidator add a validator to a model
func AddValidator(modelType reflect.Type, validators ...Validator) {
	Validators.Set(modelType, append(Validators.Get(modelType), validators...))
}

// StringValidator generic string valiator pattern
type StringValidator func(str string) bool

// NewStringValidatorFunc Create a basic string validator
func NewStringValidatorFunc(v StringValidator) FieldValidatorFunc {
	return func(value interface{}, args ...interface{}) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("Value is not a string, %v", value)
		}
		if v(str) {
			return fmt.Errorf("Value does not satisfy, %v, %v", v, str)
		}
		return nil
	}
}

// ValidEmail Email Validator
func ValidEmail(value interface{}, args ...interface{}) error {
	if _, ok := value.(string); !ok {
		return fmt.Errorf("Email not a string, %v", value)
	}
	return nil
}
