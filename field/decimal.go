package field

import (
	"database/sql/driver"
	"errors"
)

type Decimal struct {
	Decimal *Dec
	shadow  *Dec
	ShadowInit
}

func (d *Decimal) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		d.Decimal, err = NewDec(v)
		if err != nil {
			return err
		}
	case []byte:
		d.Decimal, err = NewDec(string(v))
		if err != nil {
			return err
		}
	default:
		return errors.New("error")
	}

	d.DoInit(func() {
		d.shadow = &Dec{}
		*d.shadow = *d.Decimal
	})

	return nil
}

func (d Decimal) Value() (driver.Value, error) {
	return []byte(d.Decimal.String()), nil
}

func (d Decimal) ShadowValue() (driver.Value, error) {
	return []byte(d.shadow.String()), nil
}

func (d Decimal) IsDirty() bool {
	return *d.shadow != *d.Decimal
}

func (d Decimal) IsSet() bool {
	return d.InitDone()
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	return []byte(d.Decimal.String()), nil
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	return d.Scan(data)
}
