package field

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBool(t *testing.T) {
	Convey("Scan", t, func() {
		Convey("Scan should load Bool and Shadow field", func() {
			s := &Bool{}
			s.Scan(true)

			So(s.Bool, ShouldEqual, true)
			So(s.shadow, ShouldEqual, true)
		})
		Convey("Scan should accept string", func() {
			s := Bool{}
			s.Scan("true")
			So(s.Bool, ShouldEqual, true)
		})
		Convey("Scan should accept byte", func() {
			s := Bool{}
			s.Scan([]byte("true"))
			So(s.Bool, ShouldEqual, true)
		})

		Convey("secondary Scan should not update shadow", func() {

			s := &Bool{}
			s.Scan(true)
			s.Scan(false)

			So(s.Bool, ShouldEqual, false)
			So(s.shadow, ShouldEqual, true)
		})

		Convey("Scanning sets IsSet()", func() {
			s := &Bool{}
			So(s.IsSet(), ShouldBeFalse)
			s.Scan(true)
			So(s.IsSet(), ShouldBeTrue)
		})

		Convey("empty string", func() {
			s := &NullBool{}

			err := s.Scan("")
			So(err, ShouldNotBeNil)

			v, err := s.Value()
			So(err, ShouldBeNil)
			So(v, ShouldBeNil)
		})
	})

	Convey("Value", t, func() {
		Convey("should return scanned value ", func() {
			s := &Bool{}
			s.Scan(true)

			value, err := s.Value()

			So(value, ShouldEqual, true)
			So(err, ShouldBeNil)
		})

		Convey("should return default value for type is not Scanned", func() {
			s := &Bool{}

			value, err := s.Value()
			So(value, ShouldEqual, false)
			So(err, ShouldBeNil)
		})
	})

	Convey("IsSet", t, func() {
		Convey("return false when not Scanned, true after", func() {
			s := &Bool{}
			So(s.IsSet(), ShouldBeFalse)
			s.Scan(true)
			So(s.IsSet(), ShouldBeTrue)
		})

		SkipConvey("IsSet should still be false if invalid value scanned", func() {
			s := &Bool{}
			s.Scan("pony")
			So(s.IsSet(), ShouldBeFalse)
		})

	})

	Convey("IsDirty", t, func() {
		Convey("should be false", func() {
			s := &Bool{}
			s.Scan(true)

			So(s.IsDirty(), ShouldBeFalse)
		})
		Convey("should be false if nil", func() {
			s := &Bool{}
			s.Scan(nil)

			So(s.IsDirty(), ShouldBeFalse)
		})

		Convey("should be true if modified", func() {
			s := &Bool{}
			s.Scan(true)
			s.Scan(false)

			So(s.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified from nil", func() {
			s := &Bool{}
			s.Scan(nil)   // does not set to empty string "" as it errors out.
			s.Scan(false) // sets to false

			// TODO: FIX Scan Nil Logic
			//      So(s.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified to nil", func() {
			s := &Bool{}
			s.Scan(true)
			s.Scan(nil) // nil doesn't update

			// TODO: FIX Scan Nil Logic
			//      So(s.IsDirty(), ShouldBeTrue)
		})

	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			s := &Bool{}
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
		})

		SkipConvey("should return error when only a nil scanned", func() {
			// Bool scaner does not error on invalid data
			s := &Bool{}
			s.Scan(nil)
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
		})
		Convey("should return scanned string", func() {
			s := &Bool{}
			s.Scan(true)
			value, err := s.ShadowValue()

			So(value, ShouldEqual, true)
			So(err, ShouldBeNil)
		})
	})

	Convey("MarshalJSON", t, func() {
		s := Bool{}
		s.Scan(true)
		data, err := json.Marshal(s)
		So(err, ShouldBeNil)
		So(string(data), ShouldEqual, "true")
	})

	Convey("UnmarshalJSON", t, func() {
		s := Bool{}
		err := s.UnmarshalJSON([]byte("true"))
		So(err, ShouldBeNil)
		v, err := s.Value()
		So(err, ShouldBeNil)
		So(v, ShouldEqual, true)
		So(s.IsSet(), ShouldBeTrue)
	})
}

func TestNullBool(t *testing.T) {
	Convey("Scan", t, func() {
		Convey("Scan should load Bool and Shadow field", func() {
			ns := &NullBool{}
			ns.Scan(true)

			So(ns.Bool, ShouldEqual, true)
			So(ns.shadow.Bool, ShouldEqual, true)
		})
		Convey("secondary Scan should not update shadow", func() {
			ns := &NullBool{}
			ns.Scan(true)
			ns.Scan(false)

			So(ns.Bool, ShouldEqual, false)
			So(ns.shadow.Bool, ShouldEqual, true)
		})

		Convey("IsSet set when successful", func() {
			ns := &NullBool{}
			So(ns.IsSet(), ShouldBeFalse)
			ns.Scan(true)
			So(ns.IsSet(), ShouldBeTrue)
		})
	})

	Convey("Value", t, func() {
		Convey("should return scanned value", func() {
			ns := &NullBool{}
			ns.Scan(true)

			value, err := ns.Value()

			So(value, ShouldEqual, true)
			So(err, ShouldBeNil)
		})

		Convey("should return scanned nil", func() {
			ns := &NullBool{}
			ns.Scan(nil)

			value, err := ns.Value()
			So(value, ShouldBeNil)
			So(err, ShouldBeNil)
		})

		Convey("should return nil if not set", func() {
			ns := &NullBool{}

			value, err := ns.Value()
			So(value, ShouldBeNil)
			So(err, ShouldBeNil)
		})
	})

	Convey("IsDirty", t, func() {
		Convey("should be false", func() {
			ns := &NullBool{}
			ns.Scan(true)

			So(ns.IsDirty(), ShouldBeFalse)
		})
		Convey("should be false if nil", func() {
			ns := &NullBool{}
			ns.Scan(nil)

			So(ns.IsDirty(), ShouldBeFalse)
		})

		Convey("should be true if modified", func() {
			ns := &NullBool{}
			ns.Scan(true)
			ns.Scan(false)

			So(ns.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified from nil", func() {
			ns := &NullBool{}
			ns.Scan(nil)
			ns.Scan(false)

			So(ns.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified to nil", func() {
			ns := &NullBool{}
			ns.Scan(true)
			ns.Scan(nil)

			So(ns.IsDirty(), ShouldBeTrue)
		})

	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			ns := &NullBool{}
			value, err := ns.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
		})
		Convey("should return nil before when nil", func() {
			ns := &NullBool{}
			ns.Scan(nil)
			value, err := ns.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldBeNil)
		})
		Convey("should return scanned string", func() {
			ns := &NullBool{}
			ns.Scan(true)
			value, err := ns.ShadowValue()

			So(value, ShouldEqual, true)
			So(err, ShouldBeNil)
		})
	})

	Convey("MarshalJSON", t, func() {

		Convey("when Valid", func() {
			s := NullBool{}
			s.Scan(true)
			data, err := json.Marshal(s)
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, "true")
		})

		Convey("when not Valid", func() {
			s := NullBool{}
			s.Scan(nil)
			data, err := json.Marshal(s)
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, "null")
		})

		Convey("when not set", func() {
			s := NullBool{}
			data, err := json.Marshal(s)
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, "null")
		})

	})

	Convey("UnmarshalJSON", t, func() {

		Convey("when valid bool", func() {
			s := NullBool{}
			err := json.Unmarshal([]byte("true"), &s)

			So(err, ShouldBeNil)
			v, err := s.Value()
			So(err, ShouldBeNil)
			So(v, ShouldEqual, true)
			So(s.IsSet(), ShouldBeTrue)
		})

		Convey("when null value", func() {
			s := NullBool{}
			err := json.Unmarshal([]byte("null"), &s)
			So(err, ShouldBeNil)
			_, err = s.Value()
			So(err, ShouldBeNil)
			So(s.Valid, ShouldEqual, false)
			So(s.IsSet(), ShouldBeTrue)
		})

	})
}
