package field

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// NullJSON field type, does not allow nil
type NullJSON struct {
	NullJSON interface{}
	shadow   interface{}
	isDirty  bool
	scanned  bool
	ShadowInit
}

// Scan a value into the string, error on nil
func (j *NullJSON) Scan(value interface{}) (err error) {
	switch value := value.(type) {
	case nil:
		j.NullJSON = nil
		err = nil
	case string:
		err = j.UnmarshalNullJSON([]byte(value))
	case []byte:
		err = j.UnmarshalNullJSON(value)
	case map[string]interface{}:
		_, err = json.Marshal(value)
		if err != nil {
			return err
		}
		j.NullJSON = value
		err = nil
	case []interface{}:
		_, err = json.Marshal(value)
		if err != nil {
			return err
		}
		j.NullJSON = value
		err = nil
	default:
		return errors.New("Unrecognized type")
	}

	if !j.scanned {
		j.scanned = true
	} else if !j.isDirty {
		j.isDirty = true
	}

	j.DoInit(func() {
		j.shadow = j.NullJSON
	})

	return
}

// Value return the value of this field
func (j NullJSON) Value() (driver.Value, error) {
	if j.NullJSON == nil {
		return nil, nil
	}

	bytes, err := j.MarshalNullJSON()
	if err != nil {
		return nil, err
	}

	return string(bytes), nil
}

// ShadowValue return the initial value of this field
func (j NullJSON) ShadowValue() (driver.Value, error) {
	if j.InitDone() {
		if j.shadow == nil {
			return nil, nil
		}

		bytes, err := json.Marshal(j.shadow)
		if err != nil {
			return nil, err
		}

		return string(bytes), nil
	}

	return nil, ErrorUnintializedShadow
}

// IsDirty if the shadow value does not match the field value
func (j *NullJSON) IsDirty() bool {
	return j.isDirty
}

//IsSet indicates if Scan has been called successfully
func (j NullJSON) IsSet() bool {
	return j.InitDone()
}

// MarshalNullJSON Marshal just the value of NullJSON
func (j NullJSON) MarshalNullJSON() ([]byte, error) {
	return json.Marshal(j.NullJSON)
}

// UnmarshalNullJSON implements encoding/json Unmarshaler
func (j *NullJSON) UnmarshalNullJSON(data []byte) error {
	err := json.Unmarshal(data, &j.NullJSON)
	if err != nil {
		return err
	}

	j.DoInit(func() {
		j.shadow = j.NullJSON
	})

	return nil
}
