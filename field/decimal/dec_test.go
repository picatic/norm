package decimal

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDec(t *testing.T) {
	Convey("Dec", t, func() {
		Convey("Rendering", func() {
			Convey("negative numbers should render with a minus", func() {
				d := Dec{-1234, 2}
				So(d.String(), ShouldEqual, "-12.34")
			})

			Convey("negative small decimal", func() {
				d := Dec{-1234, 5}
				So(d.String(), ShouldEqual, "-0.01234")
			})

			Convey("zero with high precision", func() {
				d := Dec{0, 5}
				So(d.String(), ShouldEqual, "0.00000")
			})
		})

		Convey("New", func() {
			Convey("valid decimal should rerender the same", func() {
				d, err := New("4.50")
				So(err, ShouldBeNil)
				So(d.Number, ShouldEqual, 450)
				So(d.Prec, ShouldEqual, 2)
			})

			Convey("negative number", func() {
				d, err := New("-2.45")
				So(err, ShouldBeNil)
				So(d.Number, ShouldEqual, -245)
				So(d.Prec, ShouldEqual, 2)
			})

			Convey("small numbers", func() {
				d, err := New("0.0245")
				So(err, ShouldBeNil)
				So(d.Number, ShouldEqual, 245)
				So(d.Prec, ShouldEqual, 4)
			})

			Convey("negative small numbers", func() {
				d, err := New("-0.000345")
				So(err, ShouldBeNil)
				So(d.Number, ShouldEqual, -345)
				So(d.Prec, ShouldEqual, 6)
			})

			Convey("invalid string should return error", func() {
				_, err := New("abc")
				So(err, ShouldNotBeNil)
			})

			Convey("multiple decimal points should return error", func() {
				_, err := New("3.4.5")
				So(err, ShouldNotBeNil)
			})

			Convey("no decimal points is valid with precision of 0", func() {
				d, err := New("450")
				So(err, ShouldBeNil)
				So(d.Number, ShouldEqual, 450)
				So(d.Prec, ShouldEqual, 0)
			})
		})

		Convey("String", func() {
			d := Dec{Number: 0, Prec: 0}
			res := d.String()
			So(res, ShouldEqual, "0")
		})

		Convey("Mul", func() {
			d := Dec{Number: 1234, Prec: 2}
			m := Dec{Number: 5678, Prec: 2}
			res := d.Mul(m).String()
			So(res, ShouldEqual, "700.6652")
		})

		Convey("Div", func() {
			Convey("Same precision", func() {
				d := Dec{Number: 1234, Prec: 2}
				m := Dec{Number: 5678, Prec: 2}
				res := d.Div(m, 2).String()
				So(res, ShouldEqual, "0.21")
			})

			Convey("Different precision", func() {
				d := Dec{Number: 1234, Prec: 3}
				m := Dec{Number: 5678, Prec: 2}
				res := d.Div(m, 3).String()
				So(res, ShouldEqual, "0.021")
			})

			Convey("Larger divide by smaller", func() {
				d := Dec{Number: 5678, Prec: 2}
				m := Dec{Number: 1234, Prec: 2}
				res := d.Div(m, 6).String()
				So(res, ShouldEqual, "4.601296")
			})
		})

		Convey("Add", func() {
			d := Dec{12345, 2}
			a := Dec{34567, 3}
			res := d.Add(a).String()
			So(res, ShouldEqual, "158.017")
		})

		Convey("Abs", func() {
			Convey("positive = positive", func() {
				d := Dec{12345, 2}
				res := d.Abs().String()
				So(res, ShouldEqual, "123.45")
			})

			Convey("negative = positive", func() {
				d := Dec{-12345, 2}
				res := d.Abs().String()
				So(res, ShouldEqual, "123.45")
			})
		})

		Convey("Neg", func() {
			Convey("positive = negative", func() {
				d := Dec{12345, 2}
				res := d.Neg().String()
				So(res, ShouldEqual, "-123.45")
			})

			Convey("negative = positive", func() {
				d := Dec{-12345, 2}
				res := d.Neg().String()
				So(res, ShouldEqual, "123.45")
			})
		})

		Convey("Sub", func() {
			Convey("positive - positive = positive", func() {
				d := Dec{12345, 2}
				s := Dec{34567, 3}
				res := d.Sub(s).String()
				So(res, ShouldEqual, "88.883")
			})

			Convey("negative - positive", func() {
				d := Dec{-12345, 2}
				s := Dec{34567, 3}
				res := d.Sub(s).String()
				So(res, ShouldEqual, "-158.017")
			})

			Convey("negative - negative", func() {
				d := Dec{-123450, 3}
				s := Dec{-34567, 2}
				res := d.Sub(s).String()
				So(res, ShouldEqual, "222.220")
			})
		})

		Convey("Lesser", func() {
			d1 := Dec{123, 2}
			d2 := Dec{234, 3}

			So(d1.Lesser(d2), ShouldBeFalse)
		})

		Convey("Greater", func() {
			d1 := Dec{123, 2}
			d2 := Dec{234, 3}

			So(d1.Greater(d2), ShouldBeTrue)
		})

		Convey("Equal", func() {
			Convey("same precision number should equal", func() {
				d1 := Dec{123, 2}
				d2 := Dec{123, 2}

				So(d1.Equals(d2), ShouldBeTrue)
			})

			Convey("27.00 should equal 27.000", func() {
				d1 := Dec{2700, 2}
				d2 := Dec{27000, 3}

				So(d1.Equals(d2), ShouldBeTrue)
			})

			Convey("27.01 should not equal 27.009", func() {
				d1 := Dec{2701, 2}
				d2 := Dec{27009, 3}

				So(d1.Equals(d2), ShouldBeFalse)
			})
		})

		Convey("Rounding", func() {
			Convey("Floor", func() {
				d := Dec{12345, 3} //12.345
				res := d.Floor(2)
				So(res.String(), ShouldEqual, "12.34")
			})

			Convey("Ceil", func() {
				d := Dec{12345, 3}
				res := d.Ceil(2)
				So(res.String(), ShouldEqual, "12.35")
			})

			Convey("Round", func() {
				Convey("Down", func() {
					d := Dec{12344, 3}
					res := d.Round(2)
					So(res.String(), ShouldEqual, "12.34")
				})

				Convey("Up", func() {
					d := Dec{12345, 3}
					res := d.Round(2)
					So(res.String(), ShouldEqual, "12.35")
				})
			})
		})
	})

	Convey("NullDec", t, func() {
		Convey("Scan", func() {
			Convey("Dec", func() {
				nd := &NullDec{}
				nd.Scan(Dec{Number: 100, Prec: 2})
				So(nd.Valid, ShouldBeTrue)
				So(nd.Dec.String(), ShouldEqual, "1.00")
			})

			Convey("*Dec", func() {
				nd := &NullDec{}
				nd.Scan(&Dec{Number: 100, Prec: 2})
				So(nd.Valid, ShouldBeTrue)
				So(nd.Dec.String(), ShouldEqual, "1.00")
			})
		})
	})
}

var _ sql.Scanner = &NullDec{}
var _ driver.Valuer = NullDec{}
