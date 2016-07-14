package field

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/picatic/norm/field/decimal"
)

type Decimal struct {
	Dec    decimal.Dec
	shadow decimal.Dec
	ShadowInit
	prec uint
}

func (d *Decimal) Scan(value interface{}) (err error) {
	value, err = ScanValuer(value)
	if err != nil {
		return err
	}

	tmp := &decimal.NullDec{}

	err = tmp.Scan(value)
	if err != nil {
		return err
	}
	if !tmp.Valid {
		ErrorCouldNotScan("Decimal", value)
	}

	d.Dec = tmp.Dec
	d.DoInit(func() {
		d.shadow = d.Dec
	})

	d.prec = d.Dec.Prec
	return nil
}

func (d Decimal) Value() (driver.Value, error) {
	dec := d.Dec.Round(d.prec)
	return dec.String(), nil
}

func (d Decimal) ShadowValue() (driver.Value, error) {
	return []byte(d.shadow.String()), nil
}

func (d Decimal) IsDirty() bool {
	return d.shadow != d.Dec
}

func (d Decimal) IsSet() bool {
	return d.InitDone()
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	fmt.Println("Marshaling:", d.Dec.String())
	return json.Marshal(d.Dec.String())
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	var numStr string
	err := json.Unmarshal(data, &numStr)
	if err != nil {
		return d.Scan(data)
	}

	return d.Scan(numStr)
}

type NullDecimal struct {
	decimal.NullDec
	shadow decimal.NullDec
	ShadowInit
}

func (d *NullDecimal) Scan(value interface{}) (err error) {
	value, err = ScanValuer(value)
	if err != nil {
		return err
	}

	err = d.NullDec.Scan(value)
	if err != nil {
		return
	}

	d.DoInit(func() {
		d.shadow = d.NullDec
	})
	return nil
}

func (d NullDecimal) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}

	return d.NullDec.Value()
}

func (d NullDecimal) ShadowValue() (driver.Value, error) {
	if !d.shadow.Valid {
		return nil, nil
	}

	return d.shadow.Value()
}

func (d NullDecimal) IsDirty() bool {
	if d.shadow.Valid && d.Valid {
		return d.shadow != d.NullDec
	}

	return d.shadow.Valid || d.Valid
}

func (d NullDecimal) IsSet() bool {
	return d.InitDone()
}

func (d NullDecimal) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(d.Dec.String())
}

func (d *NullDecimal) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return d.Scan(nil)
	}

	var numStr string
	err := json.Unmarshal(data, &numStr)
	if err != nil {
		return d.Scan(data)
	}

	return d.Scan(numStr)
}
