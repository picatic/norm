package validate

import (
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/picatic/norm/field"
	"github.com/picatic/norm/field/decimal"
)

func Nullable(validator Validator) Validator {
	return ValidatorFunc(func(v interface{}) error {
		if v == nil {
			return nil
		}

		return validator.Validate(v)
	})
}

func NotNullable(validator Validator) Validator {
	return ValidatorFunc(func(v interface{}) error {
		if v == nil {
			return NewError("property can not be nil")
		} else {
			return validator.Validate(v)
		}
	})
}

func Field(fieldName field.Name, validator Validator) Validator {
	return ValidatorFunc(func(v interface{}) error {
		value := reflect.ValueOf(v)

		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}

		if value.Kind() != reflect.Struct {
			panic("value is not a struct")
		}

		value = value.FieldByName(string(fieldName))

		if !value.IsValid() {
			panic("struct has no field " + string(fieldName))
		}

		err := validator.Validate(value.Interface())

		switch err := err.(type) {
		case ValidationError:
			err.AddLocation(fieldName)
			return err
		case ValidationErrors:
			err.AddLocation(fieldName)
			return err
		}
		return err
	})
}

//NormField validates an entry in a struct that is a Field type
func NormField(fieldName field.Name, required bool, validator Validator) Validator {
	var builder func(Validator) Validator

	if required {
		builder = NotNullable
	} else {
		builder = Nullable
	}

	return Field(fieldName, ValidatorFunc(func(v interface{}) error {
		//only need field for its driver valuer
		val, ok := v.(driver.Valuer)
		if !ok {
			panic("field is not a norm.field")
		}

		v, err := val.Value()
		if err != nil {
			panic(err)
		}

		if !required {
			switch v := v.(type) {
			case string:
				if v == "" {
					return nil
				}
			case int64:
				if v == 0 {
					return nil
				}
			}
		}

		return builder(validator).Validate(v)
	}))
}

func List(validator Validator) Validator {
	return ValidatorFunc(func(v interface{}) error {
		value := reflect.ValueOf(v)

		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}

		if value.Kind() != reflect.Slice {
			panic("property should be a slice")
		}

		var ves ValidationErrors = []*ValidationError{}
		for i := 0; i < value.Len(); i++ {
			err := validator.Validate(value.Index(i).Interface())

			if err != nil {
				switch err := err.(type) {
				case ValidationError:
					err.AddLocation(Index(i))
					ves.AddError(err)
				case ValidationErrors:
					err.AddLocation(Index(i))
					ves.AddError(err)
				}
			}
		}

		if len(ves) != 0 {
			return ves
		}

		return nil
	})
}

func IfThen(ifThis Validator, then Validator) Validator {
	return ValidatorFunc(func(v interface{}) error {
		valid := ifThis.Validate(v)

		if valid == nil {
			err := then.Validate(v)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func Length(validator Validator) Validator {
	return ValidatorFunc(func(v interface{}) (err error) {
		value := reflect.ValueOf(v)

		err = validator.Validate(value.Len())

		if err != nil {
			switch err := err.(type) {
			case ValidationError:
				err.Err = "length " + err.Err
				return err
			}

			return err
		}

		return nil
	})
}

func All(validators ...Validator) Validator {
	return ValidatorFunc(func(v interface{}) error {
		var errs ValidationErrors = []*ValidationError{}
		for _, validator := range validators {
			if err := validator.Validate(v); err != nil {
				errs.AddError(err)
			}
		}

		if len(errs) != 0 {
			return errs
		}

		return nil
	})
}

func Any(validators ...Validator) Validator {
	return ValidatorFunc(func(v interface{}) error {
		var errs ValidationErrors = []*ValidationError{}
		for _, validator := range validators {
			err := validator.Validate(v)
			if err == nil {
				return nil
			}
			// error is not nil
			errs.AddError(err)
		}

		return errs
	})
}

func First(validators ...Validator) Validator {
	return ValidatorFunc(func(v interface{}) error {
		for _, validator := range validators {
			if err := validator.Validate(v); err != nil {
				return err
			}
		}

		return nil
	})
}

func InList(list ...string) Validator {
	return ValidatorFunc(func(v interface{}) error {
		for _, item := range list {
			if v == item {
				return nil
			}
		}

		return NewError(fmt.Sprintf("%s is not in list %s", v, list))
	})
}

func NotInList(list ...string) Validator {
	return ValidatorFunc(func(v interface{}) error {
		for _, item := range list {
			if v == item {
				return NewError(fmt.Sprintf("%s is in list %s", v, list))
			}
		}

		return nil
	})
}

//Strings
func String(valName string, validator func(string) bool) Validator {
	return ValidatorFunc(func(v interface{}) error {
		str := v.(string)

		valid := validator(str)
		if !valid {
			return NewError("the string \"" + str + "\" is not a " + valName)
		}

		return nil
	})
}

//Comparisons
type comparison int

const (
	incomparable comparison = iota
	equal
	lt
	gt
)

func GT(right interface{}) Validator {
	return ValidatorFunc(func(left interface{}) error {
		c := compare(left, right)

		if c == incomparable {
			return NewError(fmt.Sprintf("%T and %T cannot be compared", left, right))
		}

		if c == gt {
			return nil
		}

		return NewError(fmt.Sprintf("property should be greater than %v", right))
	})
}

func LT(right interface{}) Validator {
	return ValidatorFunc(func(left interface{}) error {
		c := compare(left, right)

		if c == incomparable {
			return NewError(fmt.Sprintf("%T and %T cannot be compared", left, right))
		}

		if c == lt {
			return nil
		}

		return NewError(fmt.Sprintf("property should be less than %v", right))
	})
}

func GTE(right interface{}) Validator {
	return ValidatorFunc(func(left interface{}) error {
		c := compare(left, right)

		if c == incomparable {
			return NewError(fmt.Sprintf("%T and %T cannot be compared", left, right))
		}

		if c == gt || c == equal {
			return nil
		}

		return NewError(fmt.Sprintf("property should be greater or equal to %v", right))
	})
}

func LTE(right interface{}) Validator {
	return ValidatorFunc(func(left interface{}) error {
		c := compare(left, right)
		if c == incomparable {
			return NewError(fmt.Sprintf("%T and %T cannot be compared", left, right))
		}

		if c == lt || c == equal {
			return nil
		}

		return NewError(fmt.Sprintf("property should be less or equal to %v", right))
	})
}

func Equals(right interface{}) Validator {
	return ValidatorFunc(func(left interface{}) error {
		if left != right {
			return NewError(fmt.Sprintf("value should be equal to %v", right))
		}

		return nil
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
			return incomparable
		}
		if l < r {
			return lt
		} else if l > r {
			return gt
		} else if l == r {
			return equal
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
			return incomparable
		}
		if l < r {
			return lt
		} else if l > r {
			return gt
		} else if l == r {
			return equal
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
			return incomparable
		}
		if l < r {
			return lt
		} else if l > r {
			return gt
		} else if l == r {
			return equal
		}
	case reflect.String:
		l, err := decimal.New(leftValue.String())
		r, _ := decimal.New(rightValue.String())
		if err != nil {
			panic(err)
		}

		if l.Lesser(r) {
			return lt
		} else if l.Greater(r) {
			return gt
		} else if l.Equals(r) {
			return equal
		}
	}

	return incomparable
}
