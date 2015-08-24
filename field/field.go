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
func (s String) Value() (driver.Value, error) {
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
func (t *Time) Scan(value interface{}) error {
	sv, ok := value.(time.Time)
	if !ok {
		return errors.New("value should be a time and not nil")
	}

	t.DoInit(func() {
		t.shadow = sv
	})

	return nil
}

// Value return the value of this field
func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}

// ShadowValue return the initial value of this field
func (t Time) ShadowValue() (driver.Value, error) {
	if t.InitDone() {
		return t.shadow, nil
	}

	return nil, errors.New("Shadow Wasn't Created")
}

// IsDirty if the shadow value does not match the field value
func (t *Time) IsDirty() bool {
	return t.Time != t.shadow
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time)
}

// NullTime
type NullTime struct {
	dbr.NullTime
	shadow dbr.NullTime
	ShadowInit
}

// Scan a value into the Time, error on unparsable
func (nt *NullTime) Scan(value interface{}) error {

	err := nt.NullTime.Scan(value)
	if err != nil {
		return err
	}

	// load shadow on first scan only
	nt.DoInit(func() {
		_ = nt.shadow.Scan(value)
	})
	return nil
}

// Value return the value of this field
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// IsDirty if the shadow value does not match the field value
func (nt *NullTime) IsDirty() bool {
	return nt.Valid != nt.shadow.Valid || nt.Time != nt.shadow.Time
}

// ShadowValue return the initial value of this field
func (nt NullTime) ShadowValue() (driver.Value, error) {
	if nt.InitDone() {
		return nt.shadow.Value()
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
func (i *Int64) Scan(value interface{}) error {
	tmp := sql.NullInt64{}
	tmp.Scan(value)

	if tmp.Valid == false {
		// TODO: maybe nil should be simply allowed to be empty int64?
		return errors.New("value should be a int64 and not nil")
	}
	i.Int64 = tmp.Int64

	i.DoInit(func() {
		i.shadow = tmp.Int64
	})

	return nil
}

// Value return the value of this field
func (i Int64) Value() (driver.Value, error) {
	return i.Int64, nil
}

// ShadowValue return the initial value of this field
func (i Int64) ShadowValue() (driver.Value, error) {
	if i.InitDone() {
		return i.shadow, nil
	}

	return nil, errors.New("Shadow Wasn't Created")
}

// IsDirty if the shadow value does not match the field value
func (i *Int64) IsDirty() bool {
	return i.Int64 != i.shadow
}

func (i *Int64) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Int64)
}

// NullInt64
type NullInt64 struct {
	dbr.NullInt64
	shadow dbr.NullInt64
	ShadowInit
}

// Scan a value into the Int64, error on unparsable
func (ni *NullInt64) Scan(value interface{}) error {

	err := ni.NullInt64.Scan(value)
	if err != nil {
		return err
	}

	// load shadow on first scan only
	ni.DoInit(func() {
		_ = ni.shadow.Scan(value)
	})
	return nil
}

// Value return the value of this field
func (ni NullInt64) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int64, nil
}

func (ni *NullInt64) IsDirty() bool {
	return ni.Valid != ni.shadow.Valid || ni.Int64 != ni.shadow.Int64
}

// ShadowValue return the initial value of this field
func (ni NullInt64) ShadowValue() (driver.Value, error) {
	if ni.InitDone() {
		return ni.shadow.Value()
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
func (b *Bool) Scan(value interface{}) error {
	tmp := sql.NullBool{}
	tmp.Scan(value)

	if tmp.Valid == false {
		return errors.New("value should be a bool and not nil")
	}
	b.Bool = tmp.Bool

	b.DoInit(func() {
		b.shadow = tmp.Bool
	})

	return nil
}

// Value return the value of this field
func (b Bool) Value() (driver.Value, error) {
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

func (b *Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Bool)
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
