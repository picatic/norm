package field

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestScanValuer(t *testing.T) {
	Convey("ScanValuer", t, func() {
		Convey("non valuer type should pass value through", func() {
			str := "helloworld"
			str2, err := ScanValuer(str)
			So(err, ShouldBeNil)
			So(str, ShouldEqual, str2)
		})

		Convey("string field should return scanned string", func() {
			str := "helloworld"
			s := &String{}
			s.Scan(str)
			str2, err := ScanValuer(s)
			So(err, ShouldBeNil)
			So(str, ShouldEqual, str2)
		})

		Convey("Null Fields", func() {
			Convey("nil should return on nil scan", func() {
				ns := NullString{}
				ns.Scan(nil)
				n, err := ScanValuer(ns)
				So(err, ShouldBeNil)
				So(n, ShouldBeNil)
			})

			Convey("string should return on string scan", func() {
				ns := NullString{}
				ns.Scan("helloworld")
				str, err := ScanValuer(ns)
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "helloworld")
			})
		})
	})
}
