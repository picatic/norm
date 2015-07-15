package field

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/gocraft/dbr"
	"time"
)

type FieldShadow interface {
	ShadowValue() (driver.Value, error)
	IsDirty() bool
}

// FieldName, mapped to model
type FieldName string

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

type Field interface {
	sql.Scanner   // we require Scanner implementations
	driver.Valuer // our values stand and guard for thee
	FieldShadow   // we require FieldShadow
}

//
// String
//
type String struct {
	String string
	shadow string
	ShadowInit
}

func (s *String) Scan(value interface{}) error {
	sv, ok := value.(string)
	if !ok {
		// TODO: maybe nil should be simply allowed to be empty string?
		return errors.New("value should be a string and not nil")
	}
	s.String = sv

	s.DoInit(func() {
		s.shadow = sv
	})

	return nil
}

func (ns *String) Value() (driver.Value, error) {
	return ns.String, nil
}

func (ns *String) ShadowValue() (driver.Value, error) {
	if ns.InitDone() {
		return ns.shadow, nil
	}

	return nil, errors.New("Shadow Wasn't Created")
}

func (ns *String) IsDirty() bool {
	return ns.String != ns.shadow
}

//
// NullString
//
type NullString struct {
	dbr.NullString
	shadow dbr.NullString
	ShadowInit
}

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

func (ns *NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

func (ns *NullString) IsDirty() bool {
	return ns.Valid != ns.shadow.Valid || ns.String != ns.shadow.String
}

func (ns *NullString) ShadowValue() (driver.Value, error) {
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

func (ns *Time) Value() (driver.Value, error) {
	return ns.Time, nil
}

func (ns *Time) ShadowValue() (driver.Value, error) {
	if ns.InitDone() {
		return ns.shadow, nil
	}

	return nil, errors.New("Shadow Wasn't Created")
}

func (ns *Time) IsDirty() bool {
	return ns.Time != ns.shadow
}

//
// NullTime
//
type NullTime struct {
	dbr.NullTime
	shadow dbr.NullTime
	ShadowInit
}

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

func (ns *NullTime) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Time, nil
}

func (ns *NullTime) IsDirty() bool {
	return ns.Valid != ns.shadow.Valid || ns.Time != ns.shadow.Time
}

func (ns *NullTime) ShadowValue() (driver.Value, error) {
	if ns.InitDone() {
		return ns.shadow.Value()
	}
	return nil, errors.New("Shadow Wasn't Created")
}

//
// Int64
//
type Int64 struct {
	Int64  int64
	shadow int64
	ShadowInit
}

func (s *Int64) Scan(value interface{}) error {
	sv, ok := value.(int64)
	if !ok {
		// TODO: maybe nil should be simply allowed to be empty int64?
		return errors.New("value should be a int64 and not nil")
	}
	s.Int64 = sv

	s.DoInit(func() {
		s.shadow = sv
	})

	return nil
}

func (ns *Int64) Value() (driver.Value, error) {
	return ns.Int64, nil
}

func (ns *Int64) ShadowValue() (driver.Value, error) {
	if ns.InitDone() {
		return ns.shadow, nil
	}

	return nil, errors.New("Shadow Wasn't Created")
}

func (ns *Int64) IsDirty() bool {
	return ns.Int64 != ns.shadow
}

//
// NullInt64
//
type NullInt64 struct {
	dbr.NullInt64
	shadow dbr.NullInt64
	ShadowInit
}

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

func (ns *NullInt64) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Int64, nil
}

func (ns *NullInt64) IsDirty() bool {
	return ns.Valid != ns.shadow.Valid || ns.Int64 != ns.shadow.Int64
}

func (ns *NullInt64) ShadowValue() (driver.Value, error) {
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

func (s *Bool) Scan(value interface{}) error {
	sv, ok := value.(bool)
	if !ok {
		return errors.New("value should be a bool and not nil")
	}
	s.Bool = sv

	s.DoInit(func() {
		s.shadow = sv
	})

	return nil
}

func (ns *Bool) Value() (driver.Value, error) {
	return ns.Bool, nil
}

func (ns *Bool) ShadowValue() (driver.Value, error) {
	if ns.InitDone() {
		return ns.shadow, nil
	}

	return nil, errors.New("Shadow Wasn't Created")
}

func (ns *Bool) IsDirty() bool {
	return ns.Bool != ns.shadow
}

//
// NullBool
//
type NullBool struct {
	dbr.NullBool
	shadow dbr.NullBool
	ShadowInit
}

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

func (ns *NullBool) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Bool, nil
}

func (ns *NullBool) IsDirty() bool {
	return ns.Valid != ns.shadow.Valid || ns.Bool != ns.shadow.Bool
}

func (ns *NullBool) ShadowValue() (driver.Value, error) {
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
