package field

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSON field type, does not allow nil
type JSON struct {
	JSON    interface{}
	shadow  interface{}
	isDirty bool
	scanned bool
	ShadowInit
}

// Scan a value into the string, error on nil
func (j *JSON) Scan(value interface{}) (err error) {
	switch value := value.(type) {
	case nil:
		j.JSON = nil
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
		j.JSON = value
		err = nil
	case []interface{}:
		_, err = json.Marshal(value)
		if err != nil {
			return err
		}
		j.JSON = value
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
		j.shadow = j.JSON
	})

	return
}

// Value return the value of this field
func (j JSON) Value() (driver.Value, error) {
	if j.JSON == nil {
		return nil, nil
	}

	bytes, err := j.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return string(bytes), nil
}

// ShadowValue return the initial value of this field
func (j JSON) ShadowValue() (driver.Value, error) {
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
func (j *JSON) IsDirty() bool {
	return j.isDirty
}

//IsSet indicates if Scan has been called successfully
func (j JSON) IsSet() bool {
	return j.InitDone()
}

// MarshalJSON Marshal just the value of JSON
func (j JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.JSON)
}

// UnmarshalJSON implements encoding/json Unmarshaler
func (j *JSON) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &j.JSON)
	if err != nil {
		return err
	}

	j.DoInit(func() {
		j.shadow = j.JSON
	})

	return nil
}
