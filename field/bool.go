package field

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gopkg.in/guregu/null.v2"
)

// Bool that cannot be nil
type Bool struct {
	Bool   bool
	shadow bool
	Valid  bool
	ShadowInit
}

// Scan a value into the Bool, error on nil or unparsable
func (b *Bool) Scan(value interface{}) error {
	tmp := sql.NullBool{}
	tmp.Scan(value)

	if tmp.Valid == false {
		return errors.New("Value should be a bool and not nil")
	}
	b.Valid = true
	b.Bool = tmp.Bool

	b.DoInit(func() {
		b.shadow = tmp.Bool
	})

	return nil
}

// Value return the value of this field
func (b Bool) Value() (driver.Value, error) {
	if b.Valid == false {
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
func (b Bool) IsDirty() bool {
	return b.Bool != b.shadow
}

// MarshalJSON Marshal just the value of Bool
func (b Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Bool)
}

// UnmarshalJSON implements encoding/json Unmarshaler
func (b *Bool) UnmarshalJSON(data []byte) error {
	return b.Scan(string(data))
}

type nullBool null.Bool

// NullBool that can be nil
type NullBool struct {
	nullBool
	shadow null.Bool
	ShadowInit
}

// Scan a value into the Bool, error on unparsable
func (nb *NullBool) Scan(value interface{}) error {

	err := nb.nullBool.Scan(value)
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

// MarshalJSON Marshal just the value of Bool
func (nb NullBool) MarshalJSON() ([]byte, error) {
	if nb.Valid == true {
		return json.Marshal(nb.Bool)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements encoding/json Unmarshaler
func (nb *NullBool) UnmarshalJSON(data []byte) error {
	b := &null.Bool{}
	err := b.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	if b.Valid == true {
		return nb.Scan(b.Bool)
	}
	nb.Valid = false

	return nil
}
