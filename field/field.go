package field

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/gocraft/dbr"
	"time"
)

// FeildShadow Support for shadow fields. Allows us to determine if a field has been altered or not.
type FieldShadow interface {
	ShadowValue() (driver.Value, error)
	IsDirty() bool
}

// FieldName Name of a field on a model
type FieldName string

// Returns a field as SnakeCase
func (fn FieldName) SnakeCase() string {
	return dbr.NameMapping(string(fn))
}

// FieldNames
type FieldNames []FieldName

// Return []string of snake_case field names for database map
func (fn FieldNames) SnakeCase() []string {
	snakes := make([]string, len(fn))
	for i := 0; i < len(fn); i++ {
		snakes[i] = fn[i].SnakeCase()
	}
	return snakes
}

// Field Implementation required to get the basic norm features for field mapping and dirty
type Field interface {
	sql.Scanner   // we require Scanner implementations
	driver.Valuer // our values stand and guard for thee
	FieldShadow   // we require FieldShadow
}

// String field type, does not allow nil
type String struct {
	String string
	shadow string
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

	s.DoInit(func() {
		s.shadow = tmp.String
	})

	return nil
}

// Value return the value of this field
func (ns String) Value() (driver.Value, error) {
	return ns.String, nil
}

// ShadowValue return the initial value of this field
func (ns String) ShadowValue() (driver.Value, error) {
	if ns.InitDone() {
		return ns.shadow, nil
	}

	return nil, errors.New("Shadow Wasn't Created")
}

// IsDirty if the shadow value does not match the field value
func (ns *String) IsDirty() bool {
	return ns.String != ns.shadow
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

//
// Time
//
type Time struct {
	Time   time.Time
	shadow time.Time
	ShadowInit
}

// Scan a value into the Time, error on nil or unparsable
func (s *Time) Scan(value interface{}) error {
	sv, ok := value.(time.Time)
	if !ok {
		return errors.New("value should be a time and not nil")
	}

	s.DoInit(func() {
		s.shadow = sv
	})

	return nil
}

// Value return the value of this field
func (ns Time) Value() (driver.Value, error) {
	return ns.Time, nil
}

// ShadowValue return the initial value of this field
func (ns Time) ShadowValue() (driver.Value, error) {
	if ns.InitDone() {
		return ns.shadow, nil
	}

	return nil, errors.New("Shadow Wasn't Created")
}

// IsDirty if the shadow value does not match the field value
func (ns *Time) IsDirty() bool {
	return ns.Time != ns.shadow
}

func (n *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Time)
}

// NullTime
type NullTime struct {
	dbr.NullTime
	shadow dbr.NullTime
	ShadowInit
}

// Scan a value into the Time, error on unparsable
func (ns *NullTime) Scan(value interface{}) error {

	err := ns.NullTime.Scan(value)
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
func (ns NullTime) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Time, nil
}

// IsDirty if the shadow value does not match the field value
func (ns *NullTime) IsDirty() bool {
	return ns.Valid != ns.shadow.Valid || ns.Time != ns.shadow.Time
}

// ShadowValue return the initial value of this field
func (ns NullTime) ShadowValue() (driver.Value, error) {
	if ns.InitDone() {
		return ns.shadow.Value()
	}
	return nil, errors.New("Shadow Wasn't Created")
}

// Int64
type Int64 struct {
	Int64  int64
	shadow int64
	ShadowInit
}

// Scan a value into the Int64, error on nil or unparsable
func (s *Int64) Scan(value interface{}) error {
	tmp := sql.NullInt64{}
	tmp.Scan(value)

	if tmp.Valid == false {
		// TODO: maybe nil should be simply allowed to be empty int64?
		return errors.New("value should be a int64 and not nil")
	}
	s.Int64 = tmp.Int64

	s.DoInit(func() {
		s.shadow = tmp.Int64
	})

	return nil
}

// Value return the value of this field
func (ns Int64) Value() (driver.Value, error) {
	return ns.Int64, nil
}

// ShadowValue return the initial value of this field
func (ns Int64) ShadowValue() (driver.Value, error) {
	if ns.InitDone() {
		return ns.shadow, nil
	}

	return nil, errors.New("Shadow Wasn't Created")
}

// IsDirty if the shadow value does not match the field value
func (ns *Int64) IsDirty() bool {
	return ns.Int64 != ns.shadow
}

func (n *Int64) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Int64)
}

// NullInt64
type NullInt64 struct {
	dbr.NullInt64
	shadow dbr.NullInt64
	ShadowInit
}

// Scan a value into the Int64, error on unparsable
func (ns *NullInt64) Scan(value interface{}) error {

	err := ns.NullInt64.Scan(value)
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
func (ns NullInt64) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Int64, nil
}

func (ns *NullInt64) IsDirty() bool {
	return ns.Valid != ns.shadow.Valid || ns.Int64 != ns.shadow.Int64
}

// ShadowValue return the initial value of this field
func (ns NullInt64) ShadowValue() (driver.Value, error) {
	if ns.InitDone() {
		return ns.shadow.Value()
	}
	return nil, errors.New("Shadow Wasn't Created")
}

//
// Bool
//
type Bool struct {
	Bool   bool
	shadow bool
	ShadowInit
}

// Scan a value into the Bool, error on nil or unparsable
func (s *Bool) Scan(value interface{}) error {
	tmp := sql.NullBool{}
	tmp.Scan(value)

	if tmp.Valid == false {
		return errors.New("value should be a bool and not nil")
	}
	s.Bool = tmp.Bool

	s.DoInit(func() {
		s.shadow = tmp.Bool
	})

	return nil
}

// Value return the value of this field
func (ns Bool) Value() (driver.Value, error) {
	return ns.Bool, nil
}

// ShadowValue return the initial value of this field
func (ns Bool) ShadowValue() (driver.Value, error) {
	if ns.InitDone() {
		return ns.shadow, nil
	}

	return nil, errors.New("Shadow Wasn't Created")
}

// IsDirty if the shadow value does not match the field value
func (ns *Bool) IsDirty() bool {
	return ns.Bool != ns.shadow
}

func (ns *Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(ns.Bool)
}

//
// NullBool
//
type NullBool struct {
	dbr.NullBool
	shadow dbr.NullBool
	ShadowInit
}

// Scan a value into the Bool, error on unparsable
func (ns *NullBool) Scan(value interface{}) error {

	err := ns.NullBool.Scan(value)
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
func (ns NullBool) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Bool, nil
}

// IsDirty if the shadow value does not match the field value
func (ns *NullBool) IsDirty() bool {
	return ns.Valid != ns.shadow.Valid || ns.Bool != ns.shadow.Bool
}

// ShadowValue return the initial value of this field
func (ns NullBool) ShadowValue() (driver.Value, error) {
	if ns.InitDone() {
		return ns.shadow.Value()
	}
	return nil, errors.New("Shadow Wasn't Created")
}

// compile time check
var _ []sql.Scanner = []sql.Scanner{
	&String{},
	&NullString{},
	&Time{},
	&NullTime{},
	&Int64{},
	&NullInt64{},
	&Bool{},
	&NullBool{},
}

var _ []Field = []Field{
	&String{},
	&NullString{},
	&Time{},
	&NullTime{},
	&Int64{},
	&NullInt64{},
	&Bool{},
	&NullBool{},
}
