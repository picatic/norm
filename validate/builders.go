package validate

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"

	"github.com/picatic/norm"
	"github.com/picatic/norm/field"
)

func Nullable(validator Validator) Validator {
	return ValidatorFunc(func(v interface{}) error {
		if v == nil {
			return nil
		} else {
			return validator.Validate(v)
		}
	})
}

func NotNullable(validator Validator) Validator {
	return ValidatorFunc(func(v interface{}) error {
		if v == nil {
			return errors.New("value can not be nil")
		} else {
			return validator.Validate(v)
		}
	})
}

func Field(fieldName string, validator Validator) Validator {
	return ValidatorFunc(func(v interface{}) error {
		value := reflect.ValueOf(v)

		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}

		if value.Kind() != reflect.Struct {
			return errors.New("value is not a struct")
		}

		value = value.FieldByName(fieldName)

		if !value.IsValid() {
			return errors.New("struct has no field " + fieldName)
		}

		return validator.Validate(value.Interface())
	})
}

//NormField validates an entry in a struct that is a Field type
func NormField(fieldName string, validator Validator) Validator {
	return Field(fieldName, ValidatorFunc(func(v interface{}) error {
		//only need field for its driver valuer
		val, ok := v.(driver.Valuer)
		if !ok {
			return errors.New("field is not a norm.field")
		}

		v, err := val.Value()
		if err != nil {
			return err
		}

		return validator.Validate(v)
	}))
}

func List(validator Validator) Validator {
	return ValidatorFunc(func(v interface{}) error {
		value := reflect.ValueOf(v)

		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}

		if value.Kind() != reflect.Slice {
			return errors.New("value is not a slice")
		}

		for i := 0; i < value.Len(); i++ {
			err := validator.Validate(value.Index(i).Interface())
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func Length(validator Validator) Validator {
	return ValidatorFunc(func(v interface{}) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = errors.New("value does not have length")
			}
		}()

		value := reflect.ValueOf(v)

		return validator.Validate(value.Len())
	})
}

func All(validators ...Validator) Validator {
	return ValidatorFunc(func(v interface{}) error {
		for _, validator := range validators {
			if err := validator.Validate(v); err != nil {
				return err
			}
		}

		return nil
	})
}

// func And(v1 Validator, v2 Validator) Validator {
// 	return ValidatorFunc(func(v interface{}) error {
// 		e1 := v1.Validate(v)
// 		e2 := v2.Validate(v)

// 		})
// }

//NormFieldValidator this is meant as a wrapper to allow us to transition to better validation
func NormFieldValidator(fieldName field.Name, alias string, validator Validator) norm.FieldValidator {
	vFunc := func(sess norm.Session, model norm.Model, value field.Field, args ...interface{}) error {
		return validator.Validate(model)
	}

	return norm.NewFieldValidator(fieldName, alias, vFunc)
}

//Comparisons
type comparison int

const (
	equal comparison = 1 << iota
	gt
	lt
)

func GT(right interface{}) Validator {
	return ValidatorFunc(func(left interface{}) error {
		c := compare(left, right)

		if c == 0 {
			return errors.New("value is not compareable")
		}

		if c == gt {
			return nil
		}

		return fmt.Errorf("%d is not greater than %d", left, right)
	})
}

func LT(right interface{}) Validator {
	return ValidatorFunc(func(left interface{}) error {
		c := compare(left, right)

		if c == 0 {
			return errors.New("value is not compareable")
		}

		if c == lt {
			return nil
		}

		return fmt.Errorf("%d is not less than %d", left, right)
	})
}

func GTE(right interface{}) Validator {
	return ValidatorFunc(func(left interface{}) error {
		c := compare(left, right)

		if c == 0 {
			return errors.New("value is not compareable")
		}

		if c == gt || c == equal {
			return nil
		}

		return fmt.Errorf("%d is not greater than or equal to %d", left, right)
	})
}

func LTE(right interface{}) Validator {
	return ValidatorFunc(func(left interface{}) error {
		c := compare(left, right)

		if c == 0 {
			return errors.New("value is not compareable")
		}

		if c == lt || c == equal {
			return nil
		}

		return fmt.Errorf("%d is not less than or equal to %d", left, right)
	})
}

func compare(left, right interface{}) comparison {
	if left == right {
		return equal
	}

	leftValue := reflect.ValueOf(left)
	rightValue := reflect.ValueOf(right)

	switch leftValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		l := leftValue.Int()
		var r int64
		switch rightValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			r = rightValue.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			r = int64(rightValue.Uint())
		case reflect.Float32, reflect.Float64:
			r = int64(rightValue.Float())
		default:
			return comparison(0)
		}
		if l < r {
			return lt
		} else if l > r {
			return gt
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		l := leftValue.Uint()
		var r uint64
		switch rightValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			r = uint64(rightValue.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			r = rightValue.Uint()
		case reflect.Float32, reflect.Float64:
			r = uint64(rightValue.Float())
		default:
			return comparison(0)
		}
		if l < r {
			return lt
		} else if l > r {
			return gt
		}
	case reflect.Float32, reflect.Float64:
		l := leftValue.Float()
		var r float64
		switch rightValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			r = float64(rightValue.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			r = float64(rightValue.Uint())
		case reflect.Float32, reflect.Float64:
			r = rightValue.Float()
		default:
			return comparison(0)
		}
		if l < r {
			return lt
		} else if l > r {
			return gt
		}
	}

	return comparison(0)
}
