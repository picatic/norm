package field

import (
	"database/sql/driver"
	"errors"
)

type Decimal struct {
	Decimal Dec
	shadow  Dec
	ShadowInit
}

func (d *Decimal) Scan(value interface{}) (err error) {
	var dec *Dec
	switch v := value.(type) {
	case string:
		dec, err = NewDec(v)
		if err != nil {
			return err
		}
		d.Decimal = *dec
	case []byte:
		dec, err = NewDec(string(v))
		if err != nil {
			return err
		}
		d.Decimal = *dec
	default:
		return ErrorCouldNotScan("Decimal", v)
	}

	d.DoInit(func() {
		d.shadow = d.Decimal
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
	return d.shadow != d.Decimal
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

type NullDecimal struct {
	Decimal *Dec
	shadow  *Dec
	ShadowInit
}

func (d *NullDecimal) Scan(value interface{}) (err error) {
	if value != nil {
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
	} else {
		d.DoInit(func() {
			d.shadow = nil
		})
	}

	return nil
}

func (d NullDecimal) Value() (driver.Value, error) {
	if d.Decimal == nil {
		return nil, nil
	}

	return []byte(d.Decimal.String()), nil
}

func (d NullDecimal) ShadowValue() (driver.Value, error) {
	if d.shadow == nil {
		return nil, nil
	}

	return []byte(d.shadow.String()), nil
}

func (d NullDecimal) IsDirty() bool {
	if d.shadow == nil && d.Decimal == nil {
		return false
	} else if d.shadow == nil && d.Decimal != nil {
		return true
	} else if d.shadow != nil && d.Decimal == nil {
		return true
	}

	return *d.shadow != *d.Decimal
}

func (d NullDecimal) IsSet() bool {
	return d.InitDone()
}

func (d NullDecimal) MarshalJSON() ([]byte, error) {
	if d.Decimal == nil {
		return []byte("null"), nil
	}
	return []byte(d.Decimal.String()), nil
}

func (d *NullDecimal) UnmarshalJSON(data []byte) error {
	return d.Scan(data)
}
