package norm

import (
	"fmt"
	"github.com/picatic/norm/field"
	"reflect"
)

// Validators global static map of validators for models
var Validators = ValidatorMap{}

// Validator implementation
type ModelValidator interface {
	Model
	Validate(Session) error
}

type FieldValidator interface {
	Field() field.Name
	Alias() string
	Validate(Session, Model) error
}

// ValidatorMap Store validators for models
type ValidatorMap map[reflect.Type][]FieldValidator

// Get Fetch the validators registered to a Model
func (vm ValidatorMap) Get(model Model) []FieldValidator {
	modelType := reflect.TypeOf(model)
	validators, ok := vm[modelType]
	if !ok {
		validators = []FieldValidator{}
		vm[modelType] = validators
	}
	return validators
}

// Set validators for a model
func (vm ValidatorMap) Set(model Model, validators []FieldValidator) {
	modelType := reflect.TypeOf(model)
	vm[modelType] = validators
}

// Del deletes validators for a model
func (vm ValidatorMap) Del(model Model) {
	modelType := reflect.TypeOf(model)
	delete(vm, modelType)
}

// Clone the ValidatorMap
func (vm ValidatorMap) Clone() ValidatorMap {
	newMap := ValidatorMap{}
	for k, v := range vm {
		newMap[k] = make([]FieldValidator, len(v))
		for i, vv := range v {
			newMap[k][i] = vv
		}
	}
	return newMap
}

// Validate a model and specified fields. Returns nil if no errors.
func (vm ValidatorMap) Validate(sess Session, model Model, fields field.Names) *ValidationErrors {
	errs := &ValidationErrors{}
	for _, validator := range vm.Get(model) {
		switch v := validator.(type) {
		case FieldValidator:
			for _, field := range fields {
				if field == v.Field() {
					if err := v.Validate(sess, model); err != nil {
						switch err.(type) {
						case ValidationError:
							err.(*ValidationError).Field = field
							err.(*ValidationError).Alias = validator.Alias()
						default:
						}
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

// FieldValidatorFunc What a FieldValidator expects to have implemented
type FieldValidatorFunc func(sess Session, model Model, value field.Field, args ...interface{}) error

// FieldValidator Wrapper for a standard field validation
type fieldValidator struct {
	field field.Name
	alias string
	Func  FieldValidatorFunc
	Args  []interface{}
}

// NewFieldValidator Create a new FieldValidator, providing a FieldName, alias for the error,
// a function to run to validate and any args required for that validator.
func NewFieldValidator(
	field field.Name,
	alias string,
	vFunc FieldValidatorFunc,
	args ...interface{},
) FieldValidator {
	return &fieldValidator{
		field: field,
		alias: alias,
		Func:  vFunc,
		Args:  args,
	}
}

// Field that this validator is bound to
func (fv fieldValidator) Field() field.Name {
	return fv.field
}

func (fv fieldValidator) Alias() string {
	return fv.alias
}

// Validate a field on a model
func (fv fieldValidator) Validate(sess Session, model Model) error {
	field, err := ModelGetField(model, fv.Field())
	if err != nil {
		return err
	}
	return fv.Func(sess, model, field, fv.Args...)
}

// ValidationError Represent a single validation error. Contains enough information to construct a useful validation error message.
type ValidationError struct {
	Field   field.Name
	Alias   string
	Message string
}

// NewValidationError creator
func NewValidationError(f field.Name, a string, msg string) *ValidationError {
	return &ValidationError{Field: f, Alias: a, Message: msg}
}

// Error String the error
func (ve ValidationError) Error() string {
	return fmt.Sprintf("Field: [%s] Alias: [%s] Message: %s", ve.Field, ve.Alias, ve.Message)
}

// ValidationErrors Represent a set of ValidationErrors and/or error
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
// func AddValidator(modelType reflect.Type, validators ...Validator) {
// 	Validators.Set(modelType, append(Validators.Get(modelType), validators...))
// }
