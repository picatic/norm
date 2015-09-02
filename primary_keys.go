package norm

import (
	"github.com/picatic/norm/field"
)

// PrimaryKeyer a models primary key(s)
type PrimaryKeyer interface {
	Fields() field.Names
	Generator(model Model) error
}

type primaryKey struct {
	fields field.Names
}

// Fields
func (pks *primaryKey) Fields() field.Names {
	return pks.fields
}

// Generator NOOP
func (pks *primaryKey) Generator(model Model) error {
	return nil
}

// NewSinglePrimaryKey returns a single field PrimaryKeyer
func NewSinglePrimaryKey(primaryKeyField field.Name) PrimaryKeyer {
	return &primaryKey{fields: field.Names{primaryKeyField}}
}

// NewMultiplePrimaryKey returns a multiple
func NewMultiplePrimaryKey(fields field.Names) PrimaryKeyer {
	return &primaryKey{fields: fields}
}
