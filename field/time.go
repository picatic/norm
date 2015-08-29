package field

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/gocraft/dbr"
	"time"
)

// Time field that does not accept nil
type Time struct {
	Time   time.Time
	shadow time.Time
	ShadowInit
}

// Scan a value into the Time, error on nil or unparsable
func (t *Time) Scan(value interface{}) error {
	if value == nil {
		return errors.New("value cannot be nil")
	}
	tmp := dbr.NullTime{}
	tmp.Scan(value)

	if tmp.Valid == false {
		// TODO: maybe nil should be simply allowed to be empty int64?
		return errors.New("value should be a time and not nil")
	}
	t.Time = tmp.Time

	t.DoInit(func() {
		t.shadow = tmp.Time
	})

	return nil
}

// Value return the value of this field
func (t Time) Value() (driver.Value, error) {
	if t.Time.IsZero() == true {
		return nil, errors.New("Value was not set or was set to nil")
	}
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

// MarshalJSON Marshal just the value of Time
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time)
}

// NullTime time that can be nil
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
