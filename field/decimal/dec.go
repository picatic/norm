package decimal

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Dec struct {
	Number int64
	Prec   uint
}

func New(numStr string) (d *Dec, err error) {
	d = &Dec{}

	dec := strings.Split(numStr, ".")
	switch len(dec) {
	case 1:
		d.Number, err = strconv.ParseInt(dec[0], 10, 64)
		if err != nil {
			return
		}
		d.Prec = 0
		return d, nil
	case 2:
		num := strings.Join(dec, "")

		d.Number, err = strconv.ParseInt(num, 10, 64)
		if err != nil {
			return
		}
		d.Prec = uint(len(dec[1]))
		return d, nil
	default:
		return nil, errors.New("Invalid string")
	}
}

func (d Dec) Mul(mul Dec) Dec {
	return Dec{
		Number: d.Number * mul.Number,
		Prec:   d.Prec + mul.Prec,
	}
}

func (d Dec) Div(div Dec, prec uint) Dec {
	var scale int64 = 1
	scaleFact := int(prec) - int(d.Prec) + int(div.Prec)
	for i := 0; i < scaleFact; i++ {
		scale *= 10
	}
	return Dec{
		Number: d.Number * scale / div.Number,
		Prec:   prec,
	}
}

func (d Dec) Add(a Dec) Dec {
	if d.Prec < a.Prec {
		d, a = a, d
	}

	//get decimals to same precision
	for d.Prec != a.Prec {
		a.Prec++
		a.Number *= 10
	}

	return Dec{
		Number: d.Number + a.Number,
		Prec:   d.Prec,
	}
}

func (d Dec) Sub(s Dec) Dec {
	s.Number *= -1
	return d.Add(s)
}

func (d Dec) Round(prec uint) Dec {
	if d.Prec <= prec {
		return d
	}

	for d.Prec != prec {
		if d.Prec-prec == 1 {
			d.Number += 5
		}
		d.Number /= 10
		d.Prec--
	}
	return d
}

func (d Dec) Ceil(prec uint) Dec {
	if d.Prec <= prec {
		return d
	}

	for d.Prec != prec {
		if d.Prec-prec == 1 {
			d.Number += 10
		}
		d.Number /= 10
		d.Prec--
	}
	return d
}

func (d Dec) Floor(prec uint) Dec {
	if d.Prec <= prec {
		return d
	}

	for d.Prec != prec {
		d.Number /= 10
		d.Prec--
	}
	return d
}

func (d Dec) String() (str string) {
	isNegative := false
	if d.Number < 0 {
		isNegative = true
		d.Number *= -1
	}

	str = fmt.Sprintf("%d", d.Number)
	if d.Prec != 0 {
		if prec := int(d.Prec); len(str) <= prec {
			str = fmt.Sprintf("%0[1]*s", prec+1, str)
		}
		radixAt := uint(len(str)) - d.Prec
		str = str[:radixAt] + "." + str[radixAt:]
	}

	if isNegative {
		str = "-" + str
	}
	return str
}

type NullDec struct {
	Dec   Dec
	Valid bool
	Prec  uint
}

func (nd *NullDec) Scan(value interface{}) (err error) {
	var dec *Dec

	switch v := value.(type) {
	case string:
		nd.Valid = true
		dec, err = New(v)
		if err != nil {
			return err
		}
	case []byte:
		nd.Valid = true
		dec, err = New(string(v))
		if err != nil {
			return err
		}
	case nil:
		nd.Valid = false
	default:
		errors.New("could not scan")
	}

	if nd.Valid {
		nd.Dec = *dec
		nd.Prec = dec.Prec
	}
	return nil
}

func (nd NullDec) Value() (driver.Value, error) {
	if !nd.Valid {
		return nil, nil
	}

	dec := nd.Dec.Round(nd.Prec)
	return []byte(dec.String()), nil
}
