package field

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/gocraft/dbr"
)

// String field type, does not allow nil
type String struct {
	String string
	shadow string
	Valid  bool
	ShadowInit
}

// Scan a value into the string, error on nil
func (s *String) Scan(value interface{}) error {
	tmp := sql.NullString{}
	tmp.Scan(value)

	if tmp.Valid == false {
		// TODO: maybe nil should be simply allowed to be empty string?
		return errors.New("norm.field.String: value should be a string and not nil")
	}
	s.String = tmp.String
	s.Valid = true
	s.DoInit(func() {
		s.shadow = tmp.String
	})

	return nil
}

// Value return the value of this field
func (s String) Value() (driver.Value, error) {
	if s.Valid == false {
		return nil, errors.New("Value was not set or was set to nil")
	}
	return s.String, nil
}

// ShadowValue return the initial value of this field
func (s String) ShadowValue() (driver.Value, error) {
	if s.InitDone() {
		return s.shadow, nil
	}

	return nil, errors.New("Shadow Wasn't Created")
}

// IsDirty if the shadow value does not match the field value
func (s *String) IsDirty() bool {
	return s.String != s.shadow
}

// MarshalJSON Marshal just the value of String
func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String)
}

// UnmarshalJSON implements encoding/json Unmarshaler
func (s *String) UnmarshalJSON(data []byte) error {
	return s.Scan(data)
}

// NullString string that allows nil
type NullString struct {
	dbr.NullString
	shadow dbr.NullString
	ShadowInit
}

// Scan a value into the string
func (ns *NullString) Scan(value interface{}) error {

	err := ns.NullString.Scan(value)
	if err != nil {
		return err
	}

	// load shadow on first scan only
	ns.DoInit(func() {
		_ = ns.shadow.Scan(value)
	})
	return nil
}

// Value return the value of this field
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// IsDirty if the shadow value does not match the field value
func (ns *NullString) IsDirty() bool {
	return ns.Valid != ns.shadow.Valid || ns.String != ns.shadow.String
}

// ShadowValue return the initial value of this field
func (ns NullString) ShadowValue() (driver.Value, error) {
	if ns.InitDone() {
		return ns.shadow.Value()
	}
	return nil, errors.New("Shadow Wasn't Created")
}

// MarshalJSON Marshal just the value of String
func (ns NullString) MarshalJSON() ([]byte, error) {
	return json.Marshal(ns.String)
}

// UnmarshalJSON implements encoding/json Unmarshaler
func (ns *NullString) UnmarshalJSON(data []byte) error {
	return ns.Scan(data)
}
