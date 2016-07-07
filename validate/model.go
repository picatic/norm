package validate

import (
	"github.com/picatic/norm/field"
)

type ModelValidator interface {
	Validator() MapValidator
}

//MapValidator provides a way to build validators from a select list of validators
//kind of like All except you can select which fields get validated
type MapValidator map[field.Name]Validator

func (mv MapValidator) Fields(fields field.Names) Validator {
	return ValidatorFunc(func(model interface{}) error {
		for _, fieldName := range fields {
			err := mv[fieldName].Validate(model)

			if err != nil {
				return err
			}
		}

		return nil
	})
}
