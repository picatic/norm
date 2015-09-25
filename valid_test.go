package norm

import (
	"fmt"
	"github.com/picatic/norm/field"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type MockValidator struct {
	field field.Name
}

func (mv MockValidator) Field() field.Name {
	return mv.field
}

func (mv MockValidator) Alias() string {
	return "alias"
}

func (mv MockValidator) Validate(sess Session, m Model) error {
	return nil
}

func MockFieldValidatorFunc(sess Session, m Model, value field.Field, args ...interface{}) error {
	v, err := value.Value()
	if err != nil {
		return err
	}
	if v == args[0] {
		return nil
	}
	return NewValidationError("", "", fmt.Sprintf("Value [%s] did not equal first argument [%s]", v, args[0]))
}

func TestValidator(t *testing.T) {

	Convey("ValidatorMap", t, func() {
		var (
			vmap       ValidatorMap = make(ValidatorMap, 1)
			validators []FieldValidator
		)
		Convey("Get", func() {
			Convey("Not set", func() {
				So(vmap.Get(&MockModel{}), ShouldBeEmpty)
			})
			Convey("When Set", func() {
				mv := &MockValidator{}
				validators = append(validators, mv)
				vmap.Set(&MockModel{}, validators)
				So(vmap.Get(&MockModel{}), ShouldContain, mv)
			})
		})

		Convey("Set", func() {
			validators = append(validators, &MockValidator{})
			vmap.Set(&MockModel{}, validators)
			So(vmap.Get(&MockModel{}), ShouldResemble, validators)
		})

		Convey("Del", func() {
			validators = append(validators, &MockValidator{})
			vmap.Set(&MockModel{}, validators)
			vmap.Del(&MockModel{})
			So(len(vmap.Get(&MockModel{})), ShouldEqual, 0)
		})

		Convey("Clone", func() {
			validators = append(validators, &MockValidator{})
			vmap.Set(&MockModel{}, validators)
			So(vmap.Clone(), ShouldResemble, vmap)
		})

		Convey("Validate", func() {
			var (
				fv1, fv2 FieldValidator
				m        *MockModel
			)
			fv1 = NewFieldValidator(field.Name("FirstName"), "value_match", MockFieldValidatorFunc, "test")
			fv2 = NewFieldValidator(field.Name("Org"), "value_match", MockFieldValidatorFunc, "picatic")
			m = &MockModel{}
			vmap.Set(m, []FieldValidator{fv1, fv2})
			Convey("Single Error via Names", func() {
				err := vmap.Validate(nil, m, field.Names{"FirstName"})
				So(err, ShouldNotBeNil)
				So(len(err.Errors), ShouldEqual, 1)
			})

			Convey("Multiple Errors", func() {
				err := vmap.Validate(nil, m, field.Names{"FirstName", "Org"})
				So(err, ShouldNotBeNil)
				So(len(err.Errors), ShouldEqual, 2)
			})

			Convey("No errors", func() {
				m.FirstName.Scan("test")
				m.Org.Scan("picatic")
				err := vmap.Validate(nil, m, field.Names{"FirstName", "Org"})
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("ModelValidator", t, func() {

	})

	Convey("FieldValidator", t, func() {
		var (
			fv FieldValidator
			m  *MockModel
		)

		fv = NewFieldValidator(field.Name("FirstName"), "value_match", MockFieldValidatorFunc, "test")
		m = &MockModel{}
		Convey("Pass", func() {
			m.FirstName.Scan("test")
			err := fv.Validate(nil, m)
			So(err, ShouldBeNil)
		})

		Convey("Fail", func() {
			m.FirstName.Scan("duck")
			err := fv.Validate(nil, m)
			So(err, ShouldNotBeNil)
			So(fmt.Sprintf("%s", err), ShouldEqual, "Field: [] Alias: [] Message: Value [duck] did not equal first argument [test]")
		})
	})

	Convey("ValidationError", t, func() {
		var (
			ve *ValidationError
		)
		Convey("New", func() {
			ve = NewValidationError("id", "alias", "invalid")
			So(ve.Field, ShouldEqual, "id")
			So(ve.Message, ShouldEqual, "invalid")
			So(ve.Alias, ShouldEqual, "alias")
			So(ve.Error(), ShouldEqual, "Field: [id] Alias: [alias] Message: invalid")
		})
	})

	Convey("ValidationErrors", t, func() {
		var (
			ves *ValidationErrors
		)
		ves = &ValidationErrors{}
		Convey("Empty", func() {
			So(ves.Error(), ShouldEqual, "Empty errors")
		})

		Convey("Only ValidationError", func() {
			ves.Add(NewValidationError(field.Name("id"), "alias", "mega message"))
			So(ves.Error(), ShouldEqual, "First of multiple errors, Field: id Error: mega message")
		})

		Convey("Not ValidationError", func() {
			ves.Add(fmt.Errorf("not a ValidationError"))
			So(ves.Error(), ShouldEqual, "First of multiple errors is not a Validation error")
		})
	})
}
