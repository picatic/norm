package norm

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"errors"
	"reflect"
)


//var IsString ValidatorFunc = func(model ValidatableModel, field field.FieldName, args ...interface{}) <-chan error {
//	e := make(<-chan error)
//	return e
//}
var MockedError = errors.New("Mocked Error")

type MockValidator struct {
	ValidateFunc func(interface{}) error
}

func (tv *MockValidator) Validate(model interface{}) error {
	return tv.ValidateFunc(model)
}

var ErrorValidator = &MockValidator{ func(model interface{}) error {
  return MockedError
}}

var ValidValidator = &MockValidator{ func(model interface{}) error {
	return nil
}}

type ModelDouble interface {
	Name() string
}

type ModelOne struct {}

func (mo *ModelOne) Name() string { return "Model One" }

var ModelOneType = reflect.TypeOf(ModelOne{})

var ValidatorMapDouble = ValidatorMap{
	ModelOneType: []Validator{},
}

func TestValidator(t *testing.T) {
	Convey("Validator", t, func() {
		Convey("ErrorValidator", func(){
			So(ErrorValidator.Validate(nil), ShouldEqual, MockedError)
		})
		Convey("ValidValidator", func(){
			So(ValidValidator.Validate(nil), ShouldBeNil)
		})
	})
}

func TestValidatorMap(t *testing.T) {
	Convey("ValidatorMap", t, func() {
		Convey("Clone", func(){
			So(ValidatorMapDouble.Clone(), ShouldResemble, ValidatorMapDouble)
		})
	})
}
