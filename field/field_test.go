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

		Convey("Has", func() {
			fns := &Names{"Id", "FirstName", "LastName", "Created"}
			So(fns.Has(Name("Id")), ShouldBeTrue)
			So(fns.Has(Name("Email")), ShouldBeFalse)
		})

		Convey("Remove", func() {
			fns := &Names{"Id", "FirstName", "LastName", "Created"}
			So(fns.Remove(Names{"Id", "Created"}), ShouldResemble, Names{"FirstName", "LastName"})
		})
	})
}
