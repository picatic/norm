package field

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gopkg.in/guregu/null.v3"
	"time"
)

// Time field that does not accept nil
type TimeDate struct {
	Time   time.Time
	shadow time.Time
	ShadowInit
}

// Scan a value into the Time, error on nil or unparsable
func (t *TimeDate) Scan(value interface{}) error {
	var err error
	if value == nil {
		t.Time = time.Time{}
	}
	switch v := value.(type) {
	case time.Time:
		if v.IsZero() {
			return ErrorCouldNotScan("TimeDate", value)
		}
		t.Time = v
		break
	case []byte:
		t.Time, err = parseTimeDate(string(v))
		break
	case string:
		t.Time, err = parseTimeDate(v)
		break
	default:
		return ErrorCouldNotScan("TimeDate", value)
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
func (t TimeDate) Value() (driver.Value, error) {
	return t.Time, nil
}

// ShadowValue return the initial value of this field
func (t TimeDate) ShadowValue() (driver.Value, error) {
	if t.InitDone() {
		return t.shadow, nil
	}

	return nil, ErrorUnintializedShadow
}

// IsDirty if the shadow value does not match the field value
func (t *TimeDate) IsDirty() bool {
	return t.Time != t.shadow
}

//IsSet indicates if Scan has been called successfully
func (t TimeDate) IsSet() bool {
	return t.InitDone()
}

// MarshalJSON Marshal just the value of Time
func (t TimeDate) MarshalJSON() ([]byte, error) {
	str := t.Time.Format(timeDateFormat)
	return json.Marshal(str)
}

// UnmarshalJSON implements encoding/json Unmarshaler
func (t *TimeDate) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	timeDate, err := parseTimeDate(str)
	if err != nil {
		return err
	}
	return t.Scan(timeDate)
}

// NullTime time that can be nil
type NullTimeDate struct {
	nullTime
	validNull       bool
	shadow          null.Time
	shadowValidNull bool
	ShadowInit
}

// Scan a value into the Time, error on unparsable
func (nt *NullTimeDate) Scan(value interface{}) error {
	var err error
	switch v := value.(type) {
	case time.Time:
		if v.IsZero() {
			nt.Valid = false
			nt.validNull = true

		} else {
			nt.Time, nt.Valid = v, true
			nt.validNull = false
		}
		break
	case []byte:
		nt.Time, err = parseTimeDate(string(v))
		if nt.Time.IsZero() == true {
			nt.Valid = false
			nt.validNull = false
			return ErrorCouldNotScan("NullTimeDate", value)
		}
		nt.Valid = (err == nil)
		if err == nil {
			nt.validNull = false
		}
		break
	case string:
		nt.Time, err = parseTimeDate(v)
		if nt.Time.IsZero() == true {
			nt.Valid = false
			nt.validNull = false
			return ErrorCouldNotScan("NullTimeDate", value)
		}
		nt.Valid = (err == nil)
		if err == nil {
			nt.validNull = false
		}
		break
	default:
		if value == nil {
			nt.Valid = false
			nt.validNull = true
		} else {
			err = ErrorCouldNotScan("NullTimeDate", value)
		}
	}

	// load shadow on first scan only
	nt.DoInit(func() {
		_ = nt.shadow.Scan(nt.Time)
		if value == nil {
			nt.shadowValidNull = true
		}
	})
	return err
}

// Value return the value of this field
func (nt NullTimeDate) Value() (driver.Value, error) {
	if nt.validNull {
		return nil, nil
	}
	return nt.Time, nil
}

// IsDirty if the shadow value does not match the field value
func (nt *NullTimeDate) IsDirty() bool {
	if nt.validNull && nt.shadowValidNull {
		return false
	} else if nt.validNull == false && nt.shadowValidNull == false {
		return !nt.Time.Equal(nt.shadow.Time)
	}
	return true
}

//IsSet indicates if Scan has been called successfully
func (nt NullTimeDate) IsSet() bool {
	return nt.InitDone()
}

// ShadowValue return the initial value of this field
func (nt NullTimeDate) ShadowValue() (driver.Value, error) {
	if nt.InitDone() {
		if nt.shadowValidNull {
			return nil, nil
		}
		return nt.shadow.Value()
	}
	return nil, ErrorUnintializedShadow
}

// MarshalJSON Marshal just the value of String
func (nt NullTimeDate) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		str := nt.Time.Format(timeDateFormat)
		return json.Marshal(str)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements encoding/json Unmarshaler
func (nt *NullTimeDate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nt.Scan(nil)
	}

	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	timeDate, err := parseTimeDate(str)
	if err != nil {
		return err
	}
	return nt.Scan(timeDate)
}

func parseTimeDate(str string) (t time.Time, err error) {
	base := "0000-00-00"
	switch len(str) {
	case 10: // up to "YYYY-MM-DD HH:MM:SS.MMMMMM"
		if str == base[:len(str)] {
			return
		}
		t, err = time.Parse(timeDateFormat[:len(str)], str)
	default:
		err = fmt.Errorf("Invalid Time-String: %s", str)
		return
	}

	return
}
