package field

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/gocraft/dbr"
)

// Int64 that cannot be nil
type Int64 struct {
	Int64  int64
	shadow int64
	valid  bool
	ShadowInit
}

// Scan a value into the Int64, error on nil or unparsable
func (i *Int64) Scan(value interface{}) error {
	tmp := sql.NullInt64{}
	tmp.Scan(value)

	if tmp.Valid == false {
		// TODO: maybe nil should be simply allowed to be empty int64?
		return errors.New("Value should be a int64 and not nil")
	}
	i.valid = true
	i.Int64 = tmp.Int64

	i.DoInit(func() {
		i.shadow = tmp.Int64
	})

	return nil
}

// Value return the value of this field
func (i Int64) Value() (driver.Value, error) {
	if i.valid == false {
		return nil, errors.New("Value was not set or was set to nil")
	}
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

// MarshalJSON Marshal just the value of Int64
func (i Int64) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Int64)
}

// NullInt64 that can be nil
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

// IsDirty if the shadow value does not match the field value
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
