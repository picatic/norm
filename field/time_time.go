package field

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gopkg.in/guregu/null.v3"
	"time"
)

// Time field that does not accept nil
type TimeTime struct {
	Time   time.Time
	shadow time.Time
	ShadowInit
}

// Scan a value into the Time, error on nil or unparsable
func (t *TimeTime) Scan(value interface{}) error {
	value, err := ScanValuer(value)
	if err != nil {
		return err
	}

	if value == nil {
		t.Time = time.Time{}
	}
	switch v := value.(type) {
	case time.Time:
		if v.IsZero() {
			return ErrorCouldNotScan("TimeTime", value)
		}
		t.Time = v
		break
	case []byte:
		t.Time, err = parseTimeTime(string(v), time.UTC)
		break
	case string:
		t.Time, err = parseTimeTime(v, time.UTC)
		break
	default:
		return ErrorCouldNotScan("TimeTime", value)
	}
	if err != nil {
		t.Time = time.Time{}
		return err
	}
	// load shadow on first scan only
	t.DoInit(func() {
		t.shadow = t.Time
	})

	return nil
}

// Value return the value of this field
func (t TimeTime) Value() (driver.Value, error) {
	return t.Time, nil
}

// ShadowValue return the initial value of this field
func (t TimeTime) ShadowValue() (driver.Value, error) {
	if t.InitDone() {
		return t.shadow, nil
	}

	return nil, ErrorUnintializedShadow
}

// IsDirty if the shadow value does not match the field value
func (t *TimeTime) IsDirty() bool {
	return t.Time != t.shadow
}

//IsSet indicates if Scan has been called successfully
func (t TimeTime) IsSet() bool {
	return t.InitDone()
}

// MarshalJSON Marshal just the value of Time
func (t TimeTime) MarshalJSON() ([]byte, error) {
	str := t.Time.Format(timeTimeFormat)
	return json.Marshal(str)
}

// UnmarshalJSON implements encoding/json Unmarshaler
func (t *TimeTime) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	timeTime, err := parseTimeTime(str, time.UTC)
	if err != nil {
		return err
	}
	return t.Scan(timeTime)
}

// NullTime time that can be nil
type NullTimeTime struct {
	nullTime
	invalidNull       bool
	shadow            null.Time
	shadowInvalidNull bool
	ShadowInit
}

// Scan a value into the Time, error on unparsable
func (nt *NullTimeTime) Scan(value interface{}) error {
	value, err := ScanValuer(value)
	if err != nil {
		return err
	}

	switch v := value.(type) {
	case time.Time:
		nt.Time, nt.Valid = v, true
		nt.invalidNull = true
	case []byte:
		nt.Time, err = parseTimeTime(string(v), time.UTC)
		nt.Valid = (err == nil)
		if err == nil {
			nt.invalidNull = true
		}
	case string:
		nt.Time, err = parseTimeTime(v, time.UTC)
		nt.Valid = (err == nil)
		if err == nil {
			nt.invalidNull = true
		}
	default:
		if value == nil {
			nt.Valid = false
			nt.invalidNull = false
		} else {
			err = ErrorCouldNotScan("NullTimeTime", value)
		}
	}

	// load shadow on first scan only
	nt.DoInit(func() {
		_ = nt.shadow.Scan(nt.Time)
		nt.shadowInvalidNull = (value != nil)
	})
	return err
}

// Value return the value of this field
func (nt NullTimeTime) Value() (driver.Value, error) {
	if !nt.invalidNull {
		return nil, nil
	}
	return nt.Time, nil
}

// IsDirty if the shadow value does not match the field value
func (nt *NullTimeTime) IsDirty() bool {
	if !nt.invalidNull && !nt.shadowInvalidNull {
		return false
	} else if nt.invalidNull && nt.shadowInvalidNull {
		return !nt.Time.Equal(nt.shadow.Time)
	}
	return true
}

//IsSet indicates if Scan has been called successfully
func (nt NullTimeTime) IsSet() bool {
	return nt.InitDone()
}

// ShadowValue return the initial value of this field
func (nt NullTimeTime) ShadowValue() (driver.Value, error) {
	if nt.InitDone() {
		if !nt.shadowInvalidNull {
			return nil, nil
		}
		return nt.shadow.Value()
	}
	return nil, ErrorUnintializedShadow
}

// MarshalJSON Marshal just the value of String
func (nt NullTimeTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		str := nt.Time.Format(timeTimeFormat)
		return json.Marshal(str)
	}
	return json.Marshal(nil)

}

// UnmarshalJSON implements encoding/json Unmarshaler
func (nt *NullTimeTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nt.Scan(nil)
	}

	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	timeTime, err := parseTimeTime(str, time.UTC)
	if err != nil {
		return err
	}
	return nt.Scan(timeTime)
}

func parseTimeTime(str string, loc *time.Location) (t time.Time, err error) {
	base := "00:00:00"
	switch len(str) {
	case 8: // up to "YYYY-MM-DD HH:MM:SS.MMMMMM"
		if str == base[:len(str)] {
			return
		}
		t, err = time.Parse(timeTimeFormat[:len(str)], str)
	default:
		err = fmt.Errorf("Invalid Time-String: %s", str)
		return
	}

	// Adjust location
	if err == nil && loc != time.UTC {
		h, mi, s := t.Clock()
		t, err = time.Date(0, 0, 0, h, mi, s, 0, loc), nil
	}

	return
}
