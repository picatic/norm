package norm

import (
	"github.com/picatic/norm/field"
)

// PrimaryKeyer a models primary key(s)
type PrimaryKeyer interface {
	Fields() field.Names
	Generator(model Model) (field.Names, error)
}

// CustomPrimaryKeyFn function prototype for custom Generator function wrap
type CustomPrimaryKeyFn func(pk PrimaryKeyer, model Model) (field.Names, error)

func noopPrimaryKeyGenerator(pk PrimaryKeyer, model Model) (field.Names, error) {
	return field.Names{}, nil
}

type primaryKey struct {
	fields field.Names
	fn     CustomPrimaryKeyFn
}

// Fields
func (pks *primaryKey) Fields() field.Names {
	return pks.fields
}

// Generator NOOP
func (pks *primaryKey) Generator(model Model) (field.Names, error) {
	return pks.fn(pks, model)
}

// NewSinglePrimaryKey returns a single field PrimaryKeyer
func NewSinglePrimaryKey(primaryKeyField field.Name) PrimaryKeyer {
	return &primaryKey{fields: field.Names{primaryKeyField}, fn: noopPrimaryKeyGenerator}
}

// NewMultiplePrimaryKey returns a multiple
func NewMultiplePrimaryKey(fields field.Names) PrimaryKeyer {
	return &primaryKey{fields: fields, fn: noopPrimaryKeyGenerator}
}

// NewCustomPrimaryKey custom key generator
func NewCustomPrimaryKey(fields field.Names, fn CustomPrimaryKeyFn) PrimaryKeyer {
	return &primaryKey{fields: fields, fn: fn}
}
