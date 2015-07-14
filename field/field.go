package field

import (
	"database/sql"
	"errors"
	"github.com/gocraft/dbr"
	"time"
)

type FieldShadow interface {
	ShadowValue() (interface{}, error)
	IsDirty() bool
}

// FieldName, mapped to model
type FieldName string

// FieldNames
type FieldNames []FieldName

type Field interface {
	sql.Scanner // we require Scanner implementations
	FieldShadow // we require FieldShadow
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

func (ns *String) Value() (interface{}, error) {
	return ns.String, nil
}

func (ns *String) ShadowValue() (interface{}, error) {
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

func (ns *NullString) Value() (interface{}, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

func (ns *NullString) IsDirty() bool {
	return ns.Valid != ns.shadow.Valid || ns.String != ns.shadow.String
}

func (ns *NullString) ShadowValue() (interface{}, error) {
	if ns.InitDone() {
		return ns.shadow.Value()
	}
	return nil, errors.New("Shadow Wasn't Created")
}

//
// Time
//
type Time struct {
	Time time.Time
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

func (ns *Time) Value() (interface{}, error) {
	return ns.Time, nil
}

func (ns *Time) ShadowValue() (interface{}, error) {
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

func (ns *NullTime) Value() (interface{}, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Time, nil
}

func (ns *NullTime) IsDirty() bool {
	return ns.Valid != ns.shadow.Valid || ns.Time != ns.shadow.Time
}

func (ns *NullTime) ShadowValue() (interface{}, error) {
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
}

var _ []Field = []Field{
	&String{},
	&NullString{},
	&Time{},
	&NullTime{},
}

