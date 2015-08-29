package field

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/gocraft/dbr"
)

// Bool that cannot be nil
type Bool struct {
	Bool   bool
	shadow bool
	valid  bool
	ShadowInit
}

// Scan a value into the Bool, error on nil or unparsable
func (b *Bool) Scan(value interface{}) error {
	tmp := sql.NullBool{}
	tmp.Scan(value)

	if tmp.Valid == false {
		return errors.New("value should be a bool and not nil")
	}
	b.valid = true
	b.Bool = tmp.Bool

	b.DoInit(func() {
		b.shadow = tmp.Bool
	})

	return nil
}

// Value return the value of this field
func (b Bool) Value() (driver.Value, error) {
	if b.valid == false {
		return nil, errors.New("Invalid Value set")
	}
	return b.Bool, nil
}

// ShadowValue return the initial value of this field
func (b Bool) ShadowValue() (driver.Value, error) {
	if b.InitDone() {
		return b.shadow, nil
	}

	return nil, errors.New("Shadow Wasn't Created")
}

// IsDirty if the shadow value does not match the field value
func (b *Bool) IsDirty() bool {
	return b.Bool != b.shadow
}

// MarshalJSON Marshal just the value of Bool
func (b Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Bool)
}

// NullBool that can be nil
type NullBool struct {
	dbr.NullBool
	shadow dbr.NullBool
	ShadowInit
}

// Scan a value into the Bool, error on unparsable
func (nb *NullBool) Scan(value interface{}) error {

	err := nb.NullBool.Scan(value)
	if err != nil {
		return err
	}

	// load shadow on first scan only
	nb.DoInit(func() {
		_ = nb.shadow.Scan(value)
	})
	return nil
}

// Value return the value of this field
func (nb NullBool) Value() (driver.Value, error) {
	if !nb.Valid {
		return nil, nil
	}
	return nb.Bool, nil
}

// IsDirty if the shadow value does not match the field value
func (nb *NullBool) IsDirty() bool {
	return nb.Valid != nb.shadow.Valid || nb.Bool != nb.shadow.Bool
}

// ShadowValue return the initial value of this field
func (nb NullBool) ShadowValue() (driver.Value, error) {
	if nb.InitDone() {
		return nb.shadow.Value()
	}
	return nil, errors.New("Shadow Wasn't Created")
}
