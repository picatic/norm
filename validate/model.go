package validate

import (
	"github.com/picatic/norm/field"
)

type ModelValidator map[field.Name]Validator

func (mv ModelValidator) Fields(fields field.Names) Validator {
	return ValidatorFunc(func(model interface{}) error {
		for _, fieldName := range fields {
			err := Field(fieldName, mv[fieldName]).Validate(model)

			if err != nil {
				return err
			}
		}

		return nil
	})
}
