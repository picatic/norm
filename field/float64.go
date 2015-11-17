package field

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gopkg.in/guregu/null.v3"
)

// Float64 cannot be nil
type Float64 struct {
	Float64 float64
	shadow  float64
	ShadowInit
}

//Scan a value into the float64, error on nil or unparseable
func (f *Float64) Scan(value interface{}) error {
	tmp := sql.NullFloat64{}
	tmp.Scan(value)

	if tmp.Valid == false {
		return errors.New("Value should be a float64 and not nil")
	}

	f.Float64 = tmp.Float64

	f.DoInit(func() {
		f.shadow = tmp.Float64
	})

	return nil
}

//Value return the value of this field
func (f Float64) Value() (driver.Value, error) {
	return f.Float64, nil
}

//ShadowValue return the initial value of this field
func (f Float64) ShadowValue() (driver.Value, error) {
	if f.InitDone() {
		return f.shadow, nil
	}

	return nil, ErrorUnintializedShadow
}

//IsDirty if the shadow value does not match the field value
func (f *Float64) IsDirty() bool {
	return f.Float64 != f.shadow
}

//IsSet indicates if Scan has been called successfully
func (f Float64) IsSet() bool {
	return f.InitDone()
}

//MarshalJSON Marshal just the value of Int64
func (f Float64) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Float64)
}

//UnmarshalJSON implements encoding/json Unmarshaler
func (f *Float64) UnmarshalJSON(data []byte) error {
	return f.Scan(data)
}

type nullFloat null.Float

//NullFloat64 can be nil
type NullFloat64 struct {
	nullFloat
	shadow nullFloat
	ShadowInit
}

//Scan a value into the NullFloat64, error on unparseable
func (nf *NullFloat64) Scan(value interface{}) error {

	err := nf.NullFloat64.Scan(value)
	if err != nil {
		return err
	}

	//load shadow on first scan only
	nf.DoInit(func() {
		_ = nf.shadow.Scan(value)
	})
	return nil
}

//Value returns the value of the field
func (nf NullFloat64) Value() (driver.Value, error) {
	if nf.Valid == false || nf.InitDone() == false {
		return nil, nil
	}
	return nf.Float64, nil
}

//IsDirty if the shadow value does not match the field value
func (nf NullFloat64) IsDirty() bool {
	return nf.Valid != nf.shadow.Valid || nf.Float64 != nf.shadow.Float64
}

//IsSet indicates if Scan has been called successfully
func (nf NullFloat64) IsSet() bool {
	return nf.InitDone()
}

//ShadowValue returns initial value of this field value
func (nf NullFloat64) ShadowValue() (driver.Value, error) {
	if nf.InitDone() {
		return nf.shadow.Value()
	}
	return nil, ErrorUnintializedShadow
}

//MarshalJSON impliments encoding json Marshaler
func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if nf.Valid == true {
		return json.Marshal(nf.Float64)
	}
	return json.Marshal(nil)
}

//UnmarshalJSON impliments encoding/json Unmarshaler
func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	f := &null.Float{}
	err := f.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	if f.Valid == true {
		return nf.Scan(f.Float64)
	}
	return nf.Scan(nil)
}
