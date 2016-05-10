package field

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Dec struct {
	Number    int64
	Precision uint
}

func NewDec(numStr string) (d *Dec, err error) {
	d = &Dec{}

	dec := strings.Split(numStr, ".")
	switch len(dec) {
	case 1:
		d.Number, err = strconv.ParseInt(dec[0], 10, 64)
		if err != nil {
			return
		}
		d.Precision = 0
		return d, nil
	case 2:
		num := strings.Join(dec, "")

		d.Number, err = strconv.ParseInt(num, 10, 64)
		if err != nil {
			return
		}
		d.Precision = uint(len(dec[1]))
		return d, nil
	default:
		return nil, errors.New("Invalid string")
	}
}

func (d Dec) Mul(mul Dec) Dec {
	return Dec{
		Number:    d.Number * mul.Number,
		Precision: d.Precision + mul.Precision,
	}
}

func (d Dec) Div(div Dec, prec uint) Dec {
	var scale int64 = 1
	scaleFact := int(prec) - int(d.Precision) + int(div.Precision)
	for i := 0; i < scaleFact; i++ {
		scale *= 10
	}
	return Dec{
		Number:    d.Number * scale / div.Number,
		Precision: prec,
	}
}

func (d Dec) Add(a Dec) Dec {
	if d.Precision < a.Precision {
		d, a = a, d
	}

	for d.Precision != a.Precision {
		a.Precision++
		a.Number *= 10
	}

	return Dec{
		Number:    d.Number + a.Number,
		Precision: d.Precision,
	}
}

func (d Dec) Sub(s Dec) Dec {
	s.Number *= -1
	return d.Add(s)
}

func (d Dec) Round(prec uint) Dec {
	if d.Precision <= prec {
		return d
	}

	for d.Precision != prec {
		if d.Precision-prec == 1 {
			d.Number += 5
		}
		d.Number /= 10
		d.Precision--
	}
	return d
}

func (d Dec) Ceil(prec uint) Dec {
	if d.Precision <= prec {
		return d
	}

	for d.Precision != prec {
		if d.Precision-prec == 1 {
			d.Number += 10
		}
		d.Number /= 10
		d.Precision--
	}
	return d
}

func (d Dec) Floor(prec uint) Dec {
	if d.Precision <= prec {
		return d
	}

	for d.Precision != prec {
		d.Number /= 10
		d.Precision--
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
	if d.Precision == 0 {
	} else if prec := int(d.Precision); len(str) <= prec {
		str = fmt.Sprintf("%d", d.Number)
		zeros := strings.Repeat("0", prec-len(str))
		str = "0." + zeros + str
	} else {
		str = fmt.Sprintf("%d", d.Number)
		radixAt := uint(len(str)) - d.Precision
		str = str[:radixAt] + "." + str[radixAt:]
	}

	if isNegative {
		str = "-" + str
	}
	return str
}
