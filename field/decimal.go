package field

import (
	"database/sql/driver"
	"errors"
)

type Decimal struct {
	Decimal *Dec
}

func (d *Decimal) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		d.Decimal, err = NewDec(v)
		if err != nil {
			return err
		}
	default:
		return errors.New("error")
	}
	return nil
}

func (d Decimal) Value() (driver.Value, error) {
	return []byte(d.Decimal.String()), nil
}

func (d Decimal) ShadowValue() (driver.Value, error) {
	return nil, nil
}

func (d Decimal) IsDirty() bool {
	return false
}

func (d Decimal) IsSet() bool {
	return false
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	return nil
}
