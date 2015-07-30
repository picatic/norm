package norm

import (
	"errors"
	"fmt"
	"github.com/picatic/go-api/norm/field"
	"reflect"
)

type Validator interface {
	Validate(interface{}) error
}

type ValidatorMap map[reflect.Type][]Validator

func (vm ValidatorMap) Get(model interface{}) []Validator {
	modelType := reflect.TypeOf(model)
	validators, ok := vm[modelType]
	if !ok {
		validators = []Validator{}
		vm[modelType] = validators
	}
	return validators
}

func (vm ValidatorMap) Set(modelType reflect.Type, validators []Validator) {
	vm[modelType] = validators
}

func (vm ValidatorMap) Del(modelType reflect.Type) {
	delete(vm, modelType)
}

func (vm ValidatorMap) Clone() ValidatorMap {
	newMap := ValidatorMap{}
	for k, v := range vm {
		newMap[k] = []Validator{}
		copy(newMap[k], v)
	}
	return newMap
}

func (vm ValidatorMap) Validate(model interface{}, fields field.FieldNames) *ValidationErrors {
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

type FieldValidator struct {
	Field field.FieldName
	Alias string
	Func  FieldValidatorFunc
	Args  []interface{}
}

func (fv FieldValidator) Validate(model interface{}) error {
	return fv.Func(model, fv.Field, fv.Args)
}

func NewFieldValidator(
	field field.FieldName,
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

func NewGovalidator(
	field field.FieldName,
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

//func NewValidator(field field.FieldName, alias string, fn ValidatorFunc, args ...interface{}) Validator {
//	return Validator{field, alias, fn, args}
//}

//func (fv *FieldValidator) Validate(m Model, fields []field.FieldName) error {
//	return fv.Func(m, field, fv.Args...)
//}
//
//
//type ModelValidator struct {
//	Func ModelValidatorFunc
//	Args []interface{}
//}
//
//func (mv *ModelValidator) Validate(m Model, fields []field.FieldName) error {
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
//func (vl *ValidatorList) Validate(model Model, fields field.FieldNames) <-chan error {
//	err := make(chan error, 1)
//	err <- nil
//	return err
//}

// Simple
// Complex with args
// Complex Async with args
type FieldValidatorFunc func(value interface{}, args ...interface{}) error

//var vMap ValidatorMap = &ValidatorMap{"Id":{"IsString": IsString}}

type ValidationError struct {
	Field   string
	Message string
	Model   Model
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("Field: %s Error: %s", ve.Field, ve.Message)
}

type ValidationErrors struct {
	Errors []error
	Model  Model
}

func (ve *ValidationErrors) Add(err error) {
	ve.Errors = append(ve.Errors, err)
}

func (ve ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return fmt.Sprintf("Empty errors")
	}
	if fe, ok := ve.Errors[0].(*ValidationError); ok == true {
		return fmt.Sprintf("First of multiple errors, Field: %s Error: %s", fe.Field, fe.Message)
	}
	return fmt.Sprintf("First of multiple errors is not a Validation error")
}

func AddValidator(modelType reflect.Type, validators ...Validator) {
	Validators.Set(modelType, append(Validators.Get(modelType), validators...))
}

type StringValidator func(str string) bool

func NewStringValidatorFunc(v StringValidator) FieldValidatorFunc {
	return func(value interface{}, args ...interface{}) error {
		str, ok := value.(string)
		if !ok {
			return errors.New(fmt.Sprintf("Value is not a string, %v", value))
		}
		if v(str) {
			return errors.New(fmt.Sprintf("Value does not satisfy, %v, %v", v, str))
		}
		return nil
	}
}

func ValidEmail(value interface{}, args ...interface{}) error {
	if _, ok := value.(string); !ok {
		return errors.New(fmt.Sprintf("Email not a string, %v", value))
	}
	return nil
}
