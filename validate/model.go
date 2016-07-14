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
	if fields == nil {
		return ValidatorFunc(func(model interface{}) error {
			var errs ValidationErrors = []*ValidationError{}
			for _, validator := range mv {
				err := validator.Validate(model)

				if err != nil {
					errs.AddError(err)
				}
			}

			if len(errs) != 0 {
				return errs
			}

			return nil
		})
	}

	return ValidatorFunc(func(model interface{}) error {
		var errs ValidationErrors = []*ValidationError{}
		for _, fieldName := range fields {
			err := mv[fieldName].Validate(model)

			if err != nil {
				errs.AddError(err)
			}
		}

		if len(errs) != 0 {
			return errs
		}

		return nil
	})
}
