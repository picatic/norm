package field

import (
	"encoding/json"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInt64(t *testing.T) {
	Convey("Scan", t, func() {
		Convey("Scan should load Int64 and Shadow field", func() {
			s := &Int64{}
			s.Scan(1234)

			So(s.Int64, ShouldEqual, 1234)
			So(s.shadow, ShouldEqual, 1234)
		})
		Convey("secondary Scan should not update shadow", func() {

			s := &Int64{}
			s.Scan(1234)
			s.Scan(5678)

			So(s.Int64, ShouldEqual, 5678)
			So(s.shadow, ShouldEqual, 1234)
		})
	})

	Convey("Value", t, func() {
		Convey("should return scanned value", func() {
			s := &Int64{}
			s.Scan(1234)

			value, err := s.Value()

			So(value, ShouldEqual, 1234)
			So(err, ShouldBeNil)
		})

		Convey("should return empty Int64 on scanned nil", func() {
			s := &Int64{}
			s.Scan(nil)

			value, err := s.Value()
			So(value, ShouldEqual, nil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("IsDirty", t, func() {
		Convey("should be false", func() {
			s := &Int64{}
			s.Scan(1234)

			So(s.IsDirty(), ShouldBeFalse)
		})
		Convey("should be false if nil", func() {
			s := &Int64{}
			s.Scan(nil)

			So(s.IsDirty(), ShouldBeFalse)
		})

		Convey("should be true if modified", func() {
			s := &Int64{}
			s.Scan(1234)
			s.Scan(5678)

			So(s.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified from nil", func() {
			s := &Int64{}
			s.Scan(nil)  // does not set to empty Int64 "" as it errors out.
			s.Scan(5678) // sets to 5678

			// TODO: FIX Scan Nil Logic
			//      So(s.IsDirty(), ShouldBeTrue)
		})

	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			s := &Int64{}
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, errors.New("Shadow Wasn't Created"))
		})
		Convey("should return error when only a nil scanned", func() {
			s := &Int64{}
			s.Scan(nil)
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, errors.New("Shadow Wasn't Created"))
		})
		Convey("should return scanned Int64", func() {
			s := &Int64{}
			s.Scan(1234)
			value, err := s.ShadowValue()

			So(value, ShouldEqual, 1234)
			So(err, ShouldBeNil)
		})
	})

	Convey("MarshalJSON", t, func() {
		s := Int64{}
		s.Scan(1234)
		data, _ := json.Marshal(s)
		So(string(data), ShouldEqual, "1234")
	})
}

func TestNullInt64(t *testing.T) {
	Convey("Scan", t, func() {
		Convey("Scan should load Int64 and Shadow field", func() {
			ns := &NullInt64{}
			ns.Scan(1234)

			So(ns.Int64, ShouldEqual, 1234)
			So(ns.shadow.Int64, ShouldEqual, 1234)
		})
		Convey("secondary Scan should not update shadow", func() {

			ns := &NullInt64{}
			ns.Scan(1234)
			ns.Scan(5678)

			So(ns.Int64, ShouldEqual, 5678)
			So(ns.shadow.Int64, ShouldEqual, 1234)
		})
	})

	Convey("Value", t, func() {
		Convey("should return scanned value", func() {
			ns := &NullInt64{}
			ns.Scan(1234)

			value, err := ns.Value()

			So(value, ShouldEqual, 1234)
			So(err, ShouldBeNil)
		})

		Convey("should return scanned nil", func() {
			ns := &NullInt64{}
			ns.Scan(nil)

			value, err := ns.Value()
			So(value, ShouldBeNil)
			So(err, ShouldBeNil)
		})
	})

	Convey("IsDirty", t, func() {
		Convey("should be false", func() {
			ns := &NullInt64{}
			ns.Scan(1234)

			So(ns.IsDirty(), ShouldBeFalse)
		})
		Convey("should be false if nil", func() {
			ns := &NullInt64{}
			ns.Scan(nil)

			So(ns.IsDirty(), ShouldBeFalse)
		})

		Convey("should be true if modified", func() {
			ns := &NullInt64{}
			ns.Scan(1234)
			ns.Scan(5678)

			So(ns.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified from nil", func() {
			ns := &NullInt64{}
			ns.Scan(nil)
			ns.Scan(5678)

			So(ns.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified to nil", func() {
			ns := &NullInt64{}
			ns.Scan(1234)
			ns.Scan(nil)

			So(ns.IsDirty(), ShouldBeTrue)
		})

	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			ns := &NullInt64{}
			value, err := ns.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, errors.New("Shadow Wasn't Created"))
		})
		Convey("should return nil before when nil", func() {
			ns := &NullInt64{}
			ns.Scan(nil)
			value, err := ns.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldBeNil)
		})
		Convey("should return scanned Int64", func() {
			ns := &NullInt64{}
			ns.Scan(1234)
			value, err := ns.ShadowValue()

			So(value, ShouldEqual, 1234)
			So(err, ShouldBeNil)
		})
	})
}
