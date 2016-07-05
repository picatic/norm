package validate

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

type Validator interface {
	Validate(interface{}) error
}

type ValidatorFunc func(interface{}) error

func (vf ValidatorFunc) Validate(v interface{}) error {
	return vf(v)
}

var Always ValidatorFunc = func(_ interface{}) error {
	return nil
}

var Never ValidatorFunc = func(_ interface{}) error {
	return errors.New("NEVER!!!!")
}

var UUID ValidatorFunc = func(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return ErrNotString
	}

	valid := govalidator.IsUUID(str)
	if !valid {
		return errors.New("the string \"" + str + "\" is not an UUID")
	}

	return nil
}

var Email ValidatorFunc = func(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return ErrNotString
	}

	valid := govalidator.IsEmail(str)
	if !valid {
		return errors.New("the string \"" + str + "\" is not an email")
	}

	return nil
}

var URL ValidatorFunc = func(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return ErrNotString
	}

	valid := govalidator.IsURL(str)
	if !valid {
		return errors.New("the string \"" + str + "\" is not an URL")
	}
}
