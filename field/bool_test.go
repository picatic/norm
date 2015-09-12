package field

import (
	"encoding/json"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
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
	})

	Convey("Value", t, func() {
		Convey("should return scanned value", func() {
			s := &Bool{}
			s.Scan(true)

			value, err := s.Value()

			So(value, ShouldEqual, true)
			So(err, ShouldBeNil)
		})

		Convey("should return empty string on scanned nil", func() {
			s := &Bool{}
			s.Scan(nil)

			value, err := s.Value()
			So(value, ShouldEqual, nil)
			So(err, ShouldNotBeNil)
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
			So(err, ShouldResemble, errors.New("Shadow Wasn't Created"))
		})
		Convey("should return error when only a nil scanned", func() {
			s := &Bool{}
			s.Scan(nil)
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, errors.New("Shadow Wasn't Created"))
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
			So(err, ShouldResemble, errors.New("Shadow Wasn't Created"))
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
		s := NullBool{}
		s.Scan(true)
		data, err := json.Marshal(s)
		So(err, ShouldBeNil)
		So(string(data), ShouldEqual, "true")
	})

	Convey("UnmarshalJSON", t, func() {
		s := NullBool{}
		err := s.UnmarshalJSON([]byte("true"))
		So(err, ShouldBeNil)
		v, err := s.Value()
		So(err, ShouldBeNil)
		So(v, ShouldEqual, true)
	})
}
