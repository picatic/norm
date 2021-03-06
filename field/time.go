package field

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/guregu/null.v3"
	"time"
)

const (
	timeFormat     = "2006-01-02 15:04:05.000000"
	timeTimeFormat = "15:04:05"
	timeDateFormat = "2006-01-02"
)

// Time field that does not accept nil
type Time struct {
	Time   time.Time
	shadow time.Time
	ShadowInit
}

// Scan a value into the Time, error on nil or unparsable
func (t *Time) Scan(value interface{}) error {
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
			return ErrorCouldNotScan("Time", value)
		}
		t.Time = v
		break
	case []byte:
		t.Time, err = parseDateTime(string(v), time.UTC)
		break
	case string:
		t.Time, err = parseDateTime(v, time.UTC)
		break
	default:
		return ErrorCouldNotScan("Time", value)
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
func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}

// ShadowValue return the initial value of this field
func (t Time) ShadowValue() (driver.Value, error) {
	if t.InitDone() {
		return t.shadow, nil
	}

	return nil, ErrorUnintializedShadow
}

// IsDirty if the shadow value does not match the field value
func (t *Time) IsDirty() bool {
	return t.Time != t.shadow
}

//IsSet indicates if Scan has been called successfully
func (t Time) IsSet() bool {
	return t.InitDone()
}

// MarshalJSON Marshal just the value of Time
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time)
}

// UnmarshalJSON implements encoding/json Unmarshaler
func (t *Time) UnmarshalJSON(data []byte) error {
	tmp := null.Time{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	if tmp.Valid == false {
		return errors.New("Attempted to unmarshal null value")
	}
	return t.Scan(tmp.Time)
}

type nullTime null.Time

// NullTime time that can be nil
type NullTime struct {
	nullTime
	validNull       bool
	shadow          null.Time
	shadowValidNull bool
	ShadowInit
}

// Scan a value into the Time, error on unparsable
func (nt *NullTime) Scan(value interface{}) error {
	value, err := ScanValuer(value)
	if err != nil {
		return err
	}

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
		nt.Time, err = parseDateTime(string(v), time.UTC)
		if nt.Time.IsZero() == true {
			nt.Valid = false
			nt.validNull = false
			return ErrorCouldNotScan("NullTime", value)
		}
		nt.Valid = (err == nil)
		if err == nil {
			nt.validNull = false
		}
		break
	case string:
		nt.Time, err = parseDateTime(v, time.UTC)
		if nt.Time.IsZero() == true {
			nt.Valid = false
			nt.validNull = false
			return ErrorCouldNotScan("NullTime", value)
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
			err = ErrorCouldNotScan("NullTime", value)
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
func (nt NullTime) Value() (driver.Value, error) {
	if nt.validNull {
		return nil, nil
	}
	if nt.Time.IsZero() {
		return nil, nil
	}
	return nt.Time, nil
}

// IsDirty if the shadow value does not match the field value
func (nt *NullTime) IsDirty() bool {
	if nt.validNull && nt.shadowValidNull {
		return false
	} else if nt.validNull == false && nt.shadowValidNull == false {
		return !nt.Time.Equal(nt.shadow.Time)
	}
	return true
}

//IsSet indicates if Scan has been called successfully
func (nt NullTime) IsSet() bool {
	return nt.InitDone()
}

// ShadowValue return the initial value of this field
func (nt NullTime) ShadowValue() (driver.Value, error) {
	if nt.InitDone() {
		if nt.shadowValidNull {
			return nil, nil
		}
		return nt.shadow.Value()
	}
	return nil, ErrorUnintializedShadow
}

// MarshalJSON Marshal just the value of String
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.Time)
	}
	return json.Marshal(nil)

}

// UnmarshalJSON implements encoding/json Unmarshaler
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	t := &null.Time{}
	err := t.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	if t.Valid == true {
		return nt.Scan(t.Time)
	}
	return nt.Scan(nil)
}

// taken from https://github.com/go-sql-driver/mysql/blob/master/utils.go
func parseDateTime(str string, loc *time.Location) (t time.Time, err error) {
	base := "0000-00-00 00:00:00.0000000"
	switch len(str) {
	case 10, 19, 21, 22, 23, 24, 25, 26: // up to "YYYY-MM-DD HH:MM:SS.MMMMMM"
		if str == base[:len(str)] {
			return
		}
		t, err = time.Parse(timeFormat[:len(str)], str)
	default:
		err = fmt.Errorf("Invalid Time-String: %s", str)
		return
	}

	// Adjust location
	if err == nil && loc != time.UTC {
		y, mo, d := t.Date()
		h, mi, s := t.Clock()
		t, err = time.Date(y, mo, d, h, mi, s, t.Nanosecond(), loc), nil
	}

	return
}
