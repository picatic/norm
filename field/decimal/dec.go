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

func New(numStr string) (d Dec, err error) {
	d = Dec{}

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
		return Dec{}, errors.New("Invalid string")
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
	d, a = makeSamePrec(d, a)

	return Dec{
		Number: d.Number + a.Number,
		Prec:   d.Prec,
	}
}

func (d Dec) Sub(s Dec) Dec {
	s.Number *= -1
	return d.Add(s)
}

func (d Dec) Equals(e Dec) bool {
	d, e = makeSamePrec(d, e)

	return d.Number == e.Number
}

func (d Dec) Greater(gt Dec) bool {
	d, gt = makeSamePrec(d, gt)

	return d.Number > gt.Number
}

func (d Dec) Lesser(lt Dec) bool {
	d, lt = makeSamePrec(d, lt)

	return d.Number < lt.Number
}

func (d Dec) GreaterEqual(gte Dec) bool {
	d, gte = makeSamePrec(d, gte)

	return d.Number >= gte.Number
}

func (d Dec) LesserEqual(lte Dec) bool {
	d, lte = makeSamePrec(d, lte)

	return d.Number <= lte.Number
}

func makeSamePrec(d1 Dec, d2 Dec) (Dec, Dec) {
	switched := false
	if d2.Prec < d1.Prec {
		d1, d2 = d2, d1
		switched = true
	}

	//get decimals to same precision
	for d1.Prec != d2.Prec {
		d1.Prec++
		d1.Number *= 10
	}

	if switched {
		d1, d2 = d2, d1
	}

	return d1, d2
}

func (d Dec) Abs() Dec {
	if d.Number < 0 {
		d.Number *= -1
	}

	return d
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

	if d.Prec == 0 {
		str = fmt.Sprintf("%d", d.Number)
	} else {
		str = fmt.Sprintf("%0[1]*d", int(d.Prec)+1, d.Number)
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
	var dec Dec

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
	case Dec:
		nd.Valid = true
		dec = v
	case *Dec:
		nd.Valid = true
		dec = *v
	case nil:
		nd.Valid = false
	default:
		errors.New("could not scan")
	}

	if nd.Valid {
		nd.Dec = dec
		nd.Prec = dec.Prec
	}
	return nil
}

func (nd NullDec) Value() (driver.Value, error) {
	if !nd.Valid {
		return nil, nil
	}

	dec := nd.Dec.Round(nd.Prec)
	return dec.String(), nil
}
