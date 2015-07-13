package norm

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"errors"
)


// Compile time check


func TestNullString(t *testing.T) {
	Convey("Scan", t, func() {

		ns := &NullString{}
		ns.Scan("Hello")

		Convey("Scan should load String field", func() {
			So(ns.String, ShouldEqual, "Hello")
		})

		Convey("Scan should load shadow String field", func() {
			So(ns.shadow.String, ShouldEqual, "Hello")
		})
	})
	Convey("secondary Scan should not update shadow", t, func() {

		ns := &NullString{}
		ns.Scan("First")
		ns.Scan("Second")

		Convey("should update String on additional scan", func() {
			So(ns.String, ShouldEqual, "Second")
		})
		Convey("should not update shadow with additional scan", func() {
			So(ns.shadow.String, ShouldEqual, "First")
		})
	})

	Convey("Value", t, func() {

		ns := &NullString{}
		ns.Scan("First")

		Convey("should return not null value", func() {
			value, err := ns.Value()
			So(value, ShouldEqual, "First")
			So(err, ShouldBeNil)
		})
	})

	Convey("Value nil", t, func() {

		ns := &NullString{}
		ns.Scan(nil)

		Convey("should return not null value", func() {
			value, err := ns.Value()
			So(value, ShouldBeNil)
			So(err, ShouldBeNil)
		})
	})

	Convey("IsDirty", t, func() {

		Convey("should be false", func() {
			ns := &NullString{}
			ns.Scan("First")

			isDirty := ns.IsDirty()
			So(isDirty, ShouldBeFalse)
		})
		Convey("should be false if nil", func() {
			ns := &NullString{}
			ns.Scan(nil)

			isDirty := ns.IsDirty()
			So(isDirty, ShouldBeFalse)
		})

		Convey("should be true if modified", func() {
			ns := &NullString{}
			ns.Scan("First")
			ns.Scan("Second")

			isDirty := ns.IsDirty()
			So(isDirty, ShouldBeTrue)
		})
		Convey("should be true if modified from nil", func() {
			ns := &NullString{}
			ns.Scan(nil)
			ns.Scan("Second")

			isDirty := ns.IsDirty()
			So(isDirty, ShouldBeTrue)
		})
		Convey("should be true if modified to nil", func() {
			ns := &NullString{}
			ns.Scan("First")
			ns.Scan(nil)

			isDirty := ns.IsDirty()
			So(isDirty, ShouldBeTrue)
		})

	})
	Convey("ShadowValue", t, func() {


		Convey("should return err before first scan", func() {
			ns := &NullString{}
			value, err := ns.ShadowValue()
			So(value, ShouldBeNil)
			So(err, ShouldResemble, errors.New("Shadow Wasn't Created"))
		})
		Convey("should return nil before when nil", func() {
			ns := &NullString{}
			ns.Scan(nil)
			value, err := ns.ShadowValue()
			So(value, ShouldBeNil)
			So(err, ShouldBeNil)
		})
		Convey("should return scanned string", func() {
			ns := &NullString{}
			ns.Scan("First")
			value, err := ns.ShadowValue()
			So(value, ShouldEqual, "First")
			So(err, ShouldBeNil)
		})
	})

	Convey("ShadowValue nil", t, func() {

		ns := &NullString{}
		ns.Scan(nil)

		Convey("should return not null value", func() {
			value, err := ns.Value()
			So(value, ShouldBeNil)
			So(err, ShouldBeNil)
		})
	})

}