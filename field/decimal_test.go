package field

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDecimal(t *testing.T) {
	Convey("Decimal", t, func() {
		Convey("Scan", func() {
			Convey("Valid scan", func() {
				d := &Decimal{}
				err := d.Scan("4.50")
				So(err, ShouldBeNil)
				v, err := d.Value()
				So(err, ShouldBeNil)
				So(v, ShouldResemble, []byte("4.50"))
			})

			Convey("Invalid scan", func() {
				d := &Decimal{}
				err := d.Scan("klasdf")
				So(err, ShouldNotBeNil)
			})
		})

		Convey("IsSet", func() {
			d := &Decimal{}
			err := d.Scan("6.57")
			So(err, ShouldBeNil)
			So(d.IsSet(), ShouldBeTrue)
		})

		Convey("IsDirty", func() {
			Convey("Change value", func() {
				d := &Decimal{}
				d.Scan("4.50")
				d.Scan("4.55")
				So(d.IsDirty(), ShouldBeTrue)
			})

			Convey("Remain same", func() {
				d := &Decimal{}
				d.Scan("4.50")
				d.Scan("4.50")
				So(d.IsDirty(), ShouldBeFalse)
			})
		})

		Convey("MarshalJSON", func() {
			d := &Decimal{}
			d.Scan("3.40")
			bytes, err := d.MarshalJSON()
			So(err, ShouldBeNil)
			So(bytes, ShouldResemble, []byte("3.40"))
		})

		Convey("UnmarshalJSON", func() {
			d := &Decimal{}
			err := d.UnmarshalJSON([]byte("7.60"))
			So(err, ShouldBeNil)
			So(d.Decimal.Number, ShouldEqual, 760)
			So(d.Decimal.Precision, ShouldEqual, 2)
		})
	})

	Convey("NullDecimal", t, func() {
		Convey("Scanning null", func() {
			nd := &NullDecimal{}
			err := nd.Scan(nil)
			So(err, ShouldBeNil)
		})

		Convey("Marshal null", func() {
			nd := &NullDecimal{}
			nd.Scan(nil)
			bytes, err := nd.MarshalJSON()
			So(err, ShouldBeNil)
			So(bytes, ShouldResemble, []byte("null"))
		})

	})
}
