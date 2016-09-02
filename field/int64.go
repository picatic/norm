package field

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gopkg.in/guregu/null.v3"
)

// Int64 that cannot be nil
type Int64 struct {
	Int64  int64
	shadow int64
	ShadowInit
}

// Scan a value into the Int64, error on nil or unparsable
func (i *Int64) Scan(value interface{}) error {
	value, err := ScanValuer(value)
	if err != nil {
		return err
	}

	tmp := sql.NullInt64{}
	tmp.Scan(value)

	if tmp.Valid == false {
		// TODO: maybe nil should be simply allowed to be empty int64?
		return errors.New("Value should be a int64 and not nil")
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

	return nil, ErrorUnintializedShadow
}

// IsDirty if the shadow value does not match the field value
func (i *Int64) IsDirty() bool {
	return i.Int64 != i.shadow
}

//IsSet indicates if Scan has been called successfully
func (i Int64) IsSet() bool {
	return i.InitDone()
}

// MarshalJSON Marshal just the value of Int64
func (i Int64) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Int64)
}

// UnmarshalJSON implements encoding/json Unmarshaler
func (i *Int64) UnmarshalJSON(data []byte) error {
	return i.Scan(data)
}

type nullInt null.Int

// NullInt64 that can be nil
type NullInt64 struct {
	nullInt
	shadow nullInt
	ShadowInit
}

// Scan a value into the Int64, error on unparsable
func (ni *NullInt64) Scan(value interface{}) error {
	value, err := ScanValuer(value)
	if err != nil {
		return err
	}

	err = ni.NullInt64.Scan(value)
	if err != nil {
		ni.NullInt64.Valid = false
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

//IsSet indicates if Scan has been called successfully
func (ni NullInt64) IsSet() bool {
	return ni.InitDone()
}

// ShadowValue return the initial value of this field
func (ni NullInt64) ShadowValue() (driver.Value, error) {
	if ni.InitDone() {
		return ni.shadow.Value()
	}
	return nil, ErrorUnintializedShadow
}

// MarshalJSON Marshal just the value of Int64
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if ni.Valid == true {
		return json.Marshal(ni.Int64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements encoding/json Unmarshaler
func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	i := &null.Int{}
	err := i.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	if i.Valid == true {
		return ni.Scan(i.Int64)
	}
	return ni.Scan(nil)
}
