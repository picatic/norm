package field

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gopkg.in/guregu/null.v2"
)

// String field type, does not allow nil
type String struct {
	String string
	shadow string
	ShadowInit
}

// Scan a value into the string, error on nil
func (s *String) Scan(value interface{}) error {
	tmp := null.String{}
	tmp.Scan(value)

	if tmp.Valid == false {
		// TODO: maybe nil should be simply allowed to be empty string?
		return errors.New("norm.field.String: value should be a string and not nil")
	}
	s.String = tmp.String

	s.DoInit(func() {
		s.shadow = tmp.String
	})

	return nil
}

// Value return the value of this field
func (s String) Value() (driver.Value, error) {
	return s.String, nil
}

// ShadowValue return the initial value of this field
func (s String) ShadowValue() (driver.Value, error) {
	if s.InitDone() {
		return s.shadow, nil
	}

	return nil, ErrorUnintializedShadow
}

// IsDirty if the shadow value does not match the field value
func (s *String) IsDirty() bool {
	return s.String != s.shadow
}

//IsSet indicates if Scan has been called successfully
func (s String) IsSet() bool {
	return s.InitDone()
}

// MarshalJSON Marshal just the value of String
func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String)
}

// UnmarshalJSON implements encoding/json Unmarshaler
func (s *String) UnmarshalJSON(data []byte) error {
	ss := null.String{}
	err := json.Unmarshal(data, &ss)
	if err != nil {
		return nil
	}
	if ss.Valid == false {
		return errors.New("Attempted to unmarshal null value")
	}
	return s.Scan(ss.String)
}

type nullString null.String

// NullString string that allows nil
type NullString struct {
	nullString
	shadow nullString
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

//IsSet indicates if Scan has been called successfully
func (ns NullString) IsSet() bool {
	return ns.InitDone()
}

// ShadowValue return the initial value of this field
func (ns NullString) ShadowValue() (driver.Value, error) {
	if ns.InitDone() {
		return ns.shadow.Value()
	}
	return nil, ErrorUnintializedShadow
}

// MarshalJSON Marshal just the value of String
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid == true {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements encoding/json Unmarshaler
func (ns *NullString) UnmarshalJSON(data []byte) error {
	s := &null.String{}
	err := s.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	if s.Valid == true {
		return ns.Scan(s.String)
	}
	return ns.Scan(nil)

	return nil
}
