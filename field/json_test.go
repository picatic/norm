package field

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestJSON(t *testing.T) {
	Convey("JSON Scan", t, func() {
		Convey("nil", func() {
			Convey("null as string should have value nil", func() {
				js := &JSON{}
				err := js.Scan("n")
				So(err, ShouldBeNil)
				val, _ := js.Value()
				So(val, ShouldBeNil)
			})

			Convey("Scanning nil should return nil", func() {

			})
		})

		Convey("JSON Object", func() {
			js := &JSON{}
			err := js.Scan(`{"data1":1, "data2":2}`)
			So(err, ShouldBeNil)
			So(js.JSON, ShouldResemble, 2)
			So(1, ShouldEqual, 6)
		})
	})
}
