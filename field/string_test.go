package field

import (
	"encoding/json"
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
			So(s.IsSet(), ShouldBeTrue)
		})
		Convey("secondary Scan should not update shadow", func() {

			s := &String{}
			s.Scan("First")
			s.Scan("Second")

			So(s.String, ShouldEqual, "Second")
			So(s.shadow, ShouldEqual, "First")
			So(s.IsSet(), ShouldBeTrue)
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

		Convey("should return err on scanned nil", func() {
			s := &String{}
			s.Scan(nil)

			value, err := s.Value()
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "")
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

	Convey("IsSet", t, func() {
		s := &String{}
		So(s.IsSet(), ShouldBeFalse)
		s.Scan("tea pot")
		So(s.IsSet(), ShouldBeTrue)
	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			s := &String{}
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
		})
		Convey("should return error when only a nil scanned", func() {
			s := &String{}
			s.Scan(nil)
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
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

	Convey("UnmarshalJSON", t, func() {
		s := String{}
		err := json.Unmarshal([]byte("\"i am the string\""), &s)
		So(err, ShouldBeNil)
		v, err := s.Value()
		So(err, ShouldBeNil)
		So(v, ShouldEqual, "i am the string")
		So(s.IsSet(), ShouldBeTrue)
	})
}

func TestNullString(t *testing.T) {
	Convey("Scan", t, func() {
		Convey("Scan should load String and Shadow field", func() {
			ns := &NullString{}
			ns.Scan("Hello")

			So(ns.String, ShouldEqual, "Hello")
			So(ns.shadow.String, ShouldEqual, "Hello")
			So(ns.IsSet(), ShouldBeTrue)
		})
		Convey("secondary Scan should not update shadow", func() {

			ns := &NullString{}
			ns.Scan("First")
			ns.Scan("Second")

			So(ns.String, ShouldEqual, "Second")
			So(ns.shadow.String, ShouldEqual, "First")
			So(ns.IsSet(), ShouldBeTrue)
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

	Convey("IsSet", t, func() {
		s := &NullString{}
		So(s.IsSet(), ShouldBeFalse)
		s.Scan("tea pot")
		So(s.IsSet(), ShouldBeTrue)
	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			ns := &NullString{}
			value, err := ns.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
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

	Convey("MarshalJSON", t, func() {

		Convey("with valid value", func() {
			s := NullString{}
			s.Scan("Cat")
			data, _ := json.Marshal(s)
			So(string(data), ShouldEqual, "\"Cat\"")
		})

		Convey("with null value", func() {
			s := NullString{}
			data, _ := json.Marshal(s)
			So(string(data), ShouldEqual, "null")
		})
	})

	Convey("UnmarshalJSON", t, func() {

		Convey("with valid value", func() {
			s := NullString{}
			err := json.Unmarshal([]byte("\"i am the string\""), &s)
			So(err, ShouldBeNil)
			v, err := s.Value()
			So(err, ShouldBeNil)
			So(v, ShouldEqual, "i am the string")
		})

		Convey("with null value", func() {
			s := NullString{}
			err := json.Unmarshal([]byte("null"), &s)
			So(err, ShouldBeNil)
			_, err = s.Value()
			So(err, ShouldBeNil)
			So(s.Valid, ShouldBeFalse)
			So(s.IsSet(), ShouldBeTrue)
		})

	})

}
