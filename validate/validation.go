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

var UUID = String("UUID", govalidator.IsUUID)

var Email = String("email", govalidator.IsEmail)

var URL = String("URL", govalidator.IsURL)
