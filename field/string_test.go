package field

import (
	"encoding/json"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestString(t *testing.T) {
	Convey("Scan", t, func() {
		Convey("Scan should load String and Shadow field", func() {
			s := &String{}
			s.Scan("Hello")

			So(s.String, ShouldEqual, "Hello")
			So(s.shadow, ShouldEqual, "Hello")
		})
		Convey("secondary Scan should not update shadow", func() {

			s := &String{}
			s.Scan("First")
			s.Scan("Second")

			So(s.String, ShouldEqual, "Second")
			So(s.shadow, ShouldEqual, "First")
		})
	})

	Convey("Value", t, func() {
		Convey("should return scanned value", func() {
			s := &String{}
			s.Scan("First")

			value, err := s.Value()

			So(value, ShouldEqual, "First")
			So(err, ShouldBeNil)
		})

		Convey("should return empty string on scanned nil", func() {
			s := &String{}
			s.Scan(nil)

			value, err := s.Value()
			So(value, ShouldEqual, "")
			So(err, ShouldBeNil)
		})
	})

	Convey("IsDirty", t, func() {
		Convey("should be false", func() {
			s := &String{}
			s.Scan("First")

			So(s.IsDirty(), ShouldBeFalse)
		})
		Convey("should be false if nil", func() {
			s := &String{}
			s.Scan(nil)

			So(s.IsDirty(), ShouldBeFalse)
		})

		Convey("should be true if modified", func() {
			s := &String{}
			s.Scan("First")
			s.Scan("Second")

			So(s.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified from nil", func() {
			s := &String{}
			s.Scan(nil)      // does not set to empty string "" as it errors out.
			s.Scan("Second") // sets to "Second"

			// TODO: FIX Scan Nil Logic
			//      So(s.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified to nil", func() {
			s := &String{}
			s.Scan("First")
			s.Scan(nil) // nil doesn't update

			// TODO: FIX Scan Nil Logic
			//      So(s.IsDirty(), ShouldBeTrue)
		})

	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			s := &String{}
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, errors.New("Shadow Wasn't Created"))
		})
		Convey("should return error when only a nil scanned", func() {
			s := &String{}
			s.Scan(nil)
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, errors.New("Shadow Wasn't Created"))
		})
		Convey("should return scanned string", func() {
			s := &String{}
			s.Scan("First")
			value, err := s.ShadowValue()

			So(value, ShouldEqual, "First")
			So(err, ShouldBeNil)
		})
	})

	Convey("MarshalJSON", t, func() {
		s := String{}
		s.Scan("Cat")
		data, _ := json.Marshal(s)
		So(string(data), ShouldEqual, "\"Cat\"")
	})
}

func TestNullString(t *testing.T) {
	Convey("Scan", t, func() {
		Convey("Scan should load String and Shadow field", func() {
			ns := &NullString{}
			ns.Scan("Hello")

			So(ns.String, ShouldEqual, "Hello")
			So(ns.shadow.String, ShouldEqual, "Hello")
		})
		Convey("secondary Scan should not update shadow", func() {

			ns := &NullString{}
			ns.Scan("First")
			ns.Scan("Second")

			So(ns.String, ShouldEqual, "Second")
			So(ns.shadow.String, ShouldEqual, "First")
		})
	})

	Convey("Value", t, func() {
		Convey("should return scanned value", func() {
			ns := &NullString{}
			ns.Scan("First")

			value, err := ns.Value()

			So(value, ShouldEqual, "First")
			So(err, ShouldBeNil)
		})

		Convey("should return scanned nil", func() {
			ns := &NullString{}
			ns.Scan(nil)

			value, err := ns.Value()
			So(value, ShouldBeNil)
			So(err, ShouldBeNil)
		})
	})

	Convey("IsDirty", t, func() {
		Convey("should be false", func() {
			ns := &NullString{}
			ns.Scan("First")

			So(ns.IsDirty(), ShouldBeFalse)
		})
		Convey("should be false if nil", func() {
			ns := &NullString{}
			ns.Scan(nil)

			So(ns.IsDirty(), ShouldBeFalse)
		})

		Convey("should be true if modified", func() {
			ns := &NullString{}
			ns.Scan("First")
			ns.Scan("Second")

			So(ns.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified from nil", func() {
			ns := &NullString{}
			ns.Scan(nil)
			ns.Scan("Second")

			So(ns.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified to nil", func() {
			ns := &NullString{}
			ns.Scan("First")
			ns.Scan(nil)

			So(ns.IsDirty(), ShouldBeTrue)
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
}
