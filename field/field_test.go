package field

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestName(t *testing.T) {
	Convey("Name", t, func() {

		Convey("SnakeCase", func() {
			fn := Name("FirstName")
			So(fn.SnakeCase(), ShouldEqual, "first_name")
		})
	})
}

func TestNames(t *testing.T) {
	Convey("Names", t, func() {

		Convey("SnakeCase", func() {
			fns := &Names{"Id", "FirstName"}
			So(fns.SnakeCase(), ShouldResemble, []string{"id", "first_name"})
		})
	})
}
