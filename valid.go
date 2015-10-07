package norm

import (
	"bytes"
	"fmt"
	"github.com/picatic/norm/field"
	"reflect"
)

// Validators global static map of validators for models.
// To be depreciated in the near future, use the Validators interface to provide FieldValidators for models.
var Validators = ValidatorMap{}

// ModelValidators implementation for a model to define its validators
type ModelValidators interface {
	Model
	Validators() []FieldValidator
}

// Validator expects Validate(Session) to do more generic/whole model validations.
//
// It should return ValidationError or ValidationErrors on validation errors
type ModelValidator interface {
	Model
	Validate(Session) error
}

// FieldValidator defines a validation on a field.
//
// Field() provides access to the field.Name being validated
// Alias() can be set to provide a unique error alias for the error itself: CODE0001, empty_string, etc.
// Validate() is expected to return a FieldValidationError with an appropriate message set informating about the error
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
//
// Returning a ValidationErrors if any errors (including non-validation related) or nil on success
func (vm ValidatorMap) Validate(sess Session, model Model, fields field.Names) *ValidationErrors {
	errs := &ValidationErrors{}
	for _, validator := range vm.Get(model) {
		switch v := validator.(type) {
		case FieldValidator:
			for _, field := range fields {
				if field == v.Field() {
					if err := v.Validate(sess, model); err != nil {
						switch err.(type) {
						case *ValidationError:
							err.(*ValidationError).Field = field
							err.(*ValidationError).Alias = validator.Alias()
							errs.Add(err)
						case *FieldValidationError:
							errs.Add(NewValidationError(field, validator.Alias(), err.(*FieldValidationError).Message))
						default:
							errs.Add(err)
						}
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
//
// Session can be used to execute queries as part of this validation
// Model is provided if further access to other fields will be needed to validate the model
// Field (implementing Valuer) provides access to the value, but you may have to cast to work with it
// args allows you to pass configuration params to the validator: range values, array of strings to match, regex, etc.
type FieldValidatorFunc func(sess Session, model Model, value field.Field, args ...interface{}) error

// FieldValidator a private implementation for a standard field validation
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

// FieldValidationError error that a field valdiations return
type FieldValidationError struct {
	Message string
}

// NewFieldValidationError create a field validator error
func NewFieldValidationError(msg string) *FieldValidationError {
	return &FieldValidationError{Message: msg}
}

// Error returns the message
func (fve FieldValidationError) Error() string {
	return fve.Message
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
	var el = len(ve.Errors)
	if el == 0 {
		return fmt.Sprintf("Empty errors")
	}
	var out = new(bytes.Buffer)

	for i, e := range ve.Errors {
		out.WriteString(e.Error())
		if i < el-1 {
			out.WriteString("; ")
		}
	}
	return out.String()
}

// AddValidator add a validator to a model
// func AddValidator(modelType reflect.Type, validators ...Validator) {
// 	Validators.Set(modelType, append(Validators.Get(modelType), validators...))
// }
