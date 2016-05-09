package field

import (
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

		Convey("NewDec", func() {
			Convey("valid decimal should rerender the same", func() {
				d, err := NewDec("4.50")
				So(err, ShouldBeNil)
				So(d.Number, ShouldEqual, 450)
				So(d.Precision, ShouldEqual, 2)
			})

			Convey("negative number", func() {
				d, err := NewDec("-2.45")
				So(err, ShouldBeNil)
				So(d.Number, ShouldEqual, -245)
				So(d.Precision, ShouldEqual, 2)
			})

			Convey("small numbers", func() {
				d, err := NewDec("0.0245")
				So(err, ShouldBeNil)
				So(d.Number, ShouldEqual, 245)
				So(d.Precision, ShouldEqual, 4)
			})

			Convey("negative small numbers", func() {
				d, err := NewDec("-0.000345")
				So(err, ShouldBeNil)
				So(d.Number, ShouldEqual, -345)
				So(d.Precision, ShouldEqual, 6)
			})

			Convey("invalid string should return error", func() {
				_, err := NewDec("abc")
				So(err, ShouldNotBeNil)
			})

			Convey("multiple decimal points should return error", func() {
				_, err := NewDec("3.4.5")
				So(err, ShouldNotBeNil)
			})
		})

		Convey("String", func() {
			d := Dec{Number: 0, Precision: 0}
			res := d.String()
			So(res, ShouldEqual, "0")
		})

		Convey("Mul", func() {
			d := Dec{Number: 1234, Precision: 2}
			m := Dec{Number: 5678, Precision: 2}
			res := d.Mul(m).String()
			So(res, ShouldEqual, "700.6652")
		})

		Convey("Div", func() {
			Convey("Same precision", func() {
				d := Dec{Number: 1234, Precision: 2}
				m := Dec{Number: 5678, Precision: 2}
				res := d.Div(m, 2).String()
				So(res, ShouldEqual, "0.21")
			})

			Convey("Different precision", func() {
				d := Dec{Number: 1234, Precision: 3}
				m := Dec{Number: 5678, Precision: 2}
				res := d.Div(m, 3).String()
				So(res, ShouldEqual, "0.021")
			})

			Convey("Larger divide by smaller", func() {
				d := Dec{Number: 5678, Precision: 2}
				m := Dec{Number: 1234, Precision: 2}
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
}
