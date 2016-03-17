package field

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// NullJson field type, does not allow nil
type NullJson struct {
	NullJson interface{}
	shadow   interface{}
	isDirty  bool
	scanned  bool
	ShadowInit
}

// Scan a value into the string, error on nil
func (j *NullJson) Scan(value interface{}) (err error) {
	switch value := value.(type) {
	case nil:
		j.NullJson = nil
		err = nil
	case string:
		err = j.UnmarshalJSON([]byte(value))
	case []byte:
		err = j.UnmarshalJSON(value)
	case map[string]interface{}:
		_, err = json.Marshal(value)
		if err != nil {
			return err
		}
		j.NullJson = value
		err = nil
	case []interface{}:
		_, err = json.Marshal(value)
		if err != nil {
			return err
		}
		j.NullJson = value
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
		j.shadow = j.NullJson
	})

	return
}

// Value return the value of this field
func (j NullJson) Value() (driver.Value, error) {
	if j.NullJson == nil {
		return nil, nil
	}

	bytes, err := j.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return string(bytes), nil
}

// ShadowValue return the initial value of this field
func (j NullJson) ShadowValue() (driver.Value, error) {
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
func (j *NullJson) IsDirty() bool {
	return j.isDirty
}

//IsSet indicates if Scan has been called successfully
func (j NullJson) IsSet() bool {
	return j.InitDone()
}

// MarshalJSON Marshal just the value of NullJson
func (j NullJson) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.NullJson)
}

// UnmarshalJSON implements encoding/json Unmarshaler
func (j *NullJson) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &j.NullJson)
	if err != nil {
		return err
	}

	j.DoInit(func() {
		j.shadow = j.NullJson
	})

	return nil
}
