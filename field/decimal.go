package field

import (
	"database/sql/driver"

	"github.com/picatic/norm/field/decimal"
)

type Decimal struct {
	Decimal decimal.Dec
	shadow  decimal.Dec
	ShadowInit
	prec uint
}

func (d *Decimal) Scan(value interface{}) (err error) {
	tmp := &decimal.NullDec{}

	err = tmp.Scan(value)
	if err != nil {
		return err
	}
	if !tmp.Valid {
		ErrorCouldNotScan("Decimal", value)
	}

	d.Decimal = tmp.Dec
	d.DoInit(func() {
		d.shadow = d.Decimal
	})

	d.prec = d.Decimal.Prec
	return nil
}

func (d Decimal) Value() (driver.Value, error) {
	dec := d.Decimal.Round(d.prec)
	return []byte(dec.String()), nil
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
	Decimal decimal.NullDec
	shadow  decimal.NullDec
	ShadowInit
}

func (d *NullDecimal) Scan(value interface{}) (err error) {
	err = d.Decimal.Scan(value)
	if err != nil {
		return
	}

	d.DoInit(func() {
		d.shadow = d.Decimal
	})
	return nil
}

func (d NullDecimal) Value() (driver.Value, error) {
	if !d.Decimal.Valid {
		return nil, nil
	}

	return d.Decimal.Value()
}

func (d NullDecimal) ShadowValue() (driver.Value, error) {
	if !d.shadow.Valid {
		return nil, nil
	}

	return d.Decimal.Value()
}

func (d NullDecimal) IsDirty() bool {
	if !d.shadow.Valid && !d.Decimal.Valid {
		return false
	} else if !d.shadow.Valid && d.Decimal.Valid {
		return true
	} else if d.shadow.Valid && !d.Decimal.Valid {
		return true
	}

	return d.shadow != d.Decimal
}

func (d NullDecimal) IsSet() bool {
	return d.InitDone()
}

func (d NullDecimal) MarshalJSON() ([]byte, error) {
	if !d.Decimal.Valid {
		return []byte("null"), nil
	}
	return []byte(d.Decimal.Dec.String()), nil
}

func (d *NullDecimal) UnmarshalJSON(data []byte) error {
	return d.Scan(data)
}
