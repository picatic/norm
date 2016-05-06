package field

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Dec struct {
	Precision uint
	Number    int64
}

func NewDec(numStr string) (d *Dec, err error) {
	dec := strings.Split(numStr, ".")
	if len(dec) != 2 {
		return nil, errors.New("error")
	}
	num := strings.Join(dec, "")

	d = &Dec{}
	d.Number, err = strconv.ParseInt(num, 10, 64)
	if err != nil {
		return
	}
	d.Precision = uint(len(dec[1]))
	return d, nil
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
	return Dec{}
}

func (d Dec) Sub(s Dec) Dec {
	return Dec{}
}

func (d Dec) Round(prec uint) Dec {
	return Dec{}
}

func (d Dec) Ceil(prec uint) Dec {
	return Dec{}
}

func (d Dec) Floor(prec uint) Dec {
	return Dec{}
}

func (d Dec) String() string {
	str := fmt.Sprintf("%d", d.Number)
	if d.Precision == 0 {
		return str
	} else if prec := int(d.Precision); len(str) <= prec {
		zeros := strings.Repeat("0", prec-len(str))
		return "0." + zeros + str
	}
	radixAt := uint(len(str)) - d.Precision
	str = str[:radixAt] + "." + str[radixAt:]
	return str
}
