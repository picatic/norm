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

		Convey("NewNameFromSnakeCase", func() {
			So(NewNameFromSnakeCase("first_name"), ShouldEqual, "FirstName")
			So(NewNameFromSnakeCase("_meow"), ShouldEqual, "Meow")
			So(NewNameFromSnakeCase("id"), ShouldEqual, "Id")
			So(NewNameFromSnakeCase("id_"), ShouldEqual, "Id")
		})
	})
}

func TestNames(t *testing.T) {
	Convey("Names", t, func() {

		Convey("NewNamesFromString", func() {
			fns := NewNamesFromString([]string{"Id", "FirstName"})
			So(len(fns), ShouldEqual, 2)
			So(fns.Has("Id"), ShouldBeTrue)
			So(fns.Has("FirstName"), ShouldBeTrue)
		})

		Convey("NewNamesFromSnakeCase", func() {
			fns := NewNamesFromSnakeCase([]string{"id", "first_name", "a_good_idea"})
			So(len(fns), ShouldEqual, 3)
			So(fns.Has("Id"), ShouldBeTrue)
			So(fns.Has("FirstName"), ShouldBeTrue)
			So(fns.Has("AGoodIdea"), ShouldBeTrue)
		})

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
			fns := Names{"Id", "FirstName", "LastName", "Created"}
			So(fns.Remove(Names{"Id", "Created"}), ShouldResemble, Names{"FirstName", "LastName"})
			So(fns.Remove(Names{}), ShouldResemble, fns)
		})

		Convey("Add", func() {
			fns := Names{"FirstName", "LastName", "Created"}
			So(fns.Add(Names{"Id"}), ShouldResemble, Names{"FirstName", "LastName", "Created", "Id"})
			So(fns.Add(Names{}), ShouldResemble, fns)
		})

		Convey("Intersect", func() {
			fns1 := Names{"Id", "FirstName", "LastName", "Created", "Modified"}
			fns2 := Names{"FirstName", "LastName", "Company"}
			fnsr := fns1.Intersect(fns2)
			So(len(fnsr), ShouldEqual, 2)
			So(fnsr, ShouldContain, "FirstName")
			So(fnsr, ShouldContain, "LastName")

			// and backwords
			fnsr = fns2.Intersect(fns1)
			So(len(fnsr), ShouldEqual, 2)
			So(fnsr, ShouldContain, "FirstName")
			So(fnsr, ShouldContain, "LastName")
		})

	})
}
