package field

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDecimal(t *testing.T) {
	Convey("Decimal", t, func() {
		Convey("Scan", func() {
			d := &Decimal{}
			d.Scan("4.50")
			v, err := d.Value()
			So(err, ShouldBeNil)
			So(v, ShouldResemble, []byte("4.50"))
		})
	})
}
