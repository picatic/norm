package valid

import "fmt"
import (
	"github.com/picatic/go-api/norm"
	"github.com/picatic/go-api/norm/field"
	"sync"
)

type ValidatorList []Validator

type Validator struct {
	Field field.FieldName
	Alias string
	Func  ValidatorFunc
	Args  []interface{}
	//	Required bool
	//	Order int32
	//	Level string
}

func (v Validator) String() string {
	return fmt.Sprintf("%s %s(%s)", v.Alias, "func", v.Args)
}

func NewValidator(field field.FieldName, alias string, fn ValidatorFunc, args ...interface{}) Validator {
	return Validator{field, alias, fn, args}
}

// Validatable
type Validatable interface {
	ValidatorList() ValidatorList // Must be locking
	ValidatorListReset()          // Reset ValidatorMap
}

type ValidatableModel interface {
	Validatable
	norm.Model
}

type ValidatorList []Validator

// Add a validator
func (vl *ValidatorList) AddValidator(field field.FieldName, alias string, fn ValidatorFunc, args ...interface{}) {
	vl = append(vl, NewValidator(field, alias, fn, args))
}

// Validate a model with this map for all fields indicated. Nil returned if no validator errors
func (vl *ValidatorList) Validate(model norm.Model, fields field.FieldNames) <-chan error {
	err := make(chan error, 1)
	err <- nil
	return err
}

// Simple
// Complex with args
// Complex Async with args
type ValidatorFunc func(model ValidatableModel, field field.FieldName, args ...interface{}) <-chan error

func (vf *ValidatorFunc) Validate() {
	errChan := vf()
	err := <-errChan
	return err != nil, err
}

var IsString ValidatorFunc = func(model ValidatableModel, field field.FieldName, args ...interface{}) <-chan error {
	e := make(<-chan error)
	return e
}

//var vMap ValidatorMap = &ValidatorMap{"Id":{"IsString": IsString}}

type ValidationError struct {
	Field   string
	Message string
	Model   norm.Model
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("Field: %s Error: %s", ve.Field, ve.Message)
}

type ValidationModelError struct {
	Message   string
	Model     norm.Model
	Validator Validator
}

func (vme *ValidationModelError) Error() string {
	return fmt.Sprintf("Model %+v Validator: %s", vme.Model, vme.Validator.String())
}

type ValidationErrors struct {
	Errors []error
	Model  norm.Model
}

func (ve *ValidationErrors) Error() string {
	return "NOPE"
}
