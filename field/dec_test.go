package field

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDec(t *testing.T) {
	Convey("Dec", t, func() {
		Convey("NewDec", func() {
			Convey("valid decimal should rerender the same", func() {
				d, err := NewDec("4.50")
				So(err, ShouldBeNil)
				So(d.String(), ShouldEqual, "4.50")
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
	})
}
