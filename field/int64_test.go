package field

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInt64(t *testing.T) {
	Convey("Scan", t, func() {
		Convey("Scan should load Int64 and Shadow field", func() {
			s := &Int64{}
			s.Scan(1234)

			So(s.Int64, ShouldEqual, 1234)
			So(s.shadow, ShouldEqual, 1234)
			So(s.IsSet(), ShouldBeTrue)
		})
		Convey("secondary Scan should not update shadow", func() {

			s := &Int64{}
			s.Scan(1234)
			s.Scan(5678)

			So(s.Int64, ShouldEqual, 5678)
			So(s.shadow, ShouldEqual, 1234)
			So(s.IsSet(), ShouldBeTrue)
		})

		Convey("empty string", func() {
			s := &NullInt64{}

			err := s.Scan("")
			So(err, ShouldNotBeNil)

			v, err := s.Value()
			So(err, ShouldBeNil)
			So(v, ShouldBeNil)
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

		Convey("should return default Int64 on scanned nil", func() {
			s := &Int64{}
			s.Scan(nil)

			value, err := s.Value()
			So(value, ShouldEqual, int64(0))
			So(err, ShouldBeNil)
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

	Convey("IsSet", t, func() {
		s := &Int64{}
		So(s.IsSet(), ShouldBeFalse)
		s.Scan(123)
		So(s.IsSet(), ShouldBeTrue)
	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			s := &Int64{}
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
		})
		Convey("should return error when only a nil scanned", func() {
			s := &Int64{}
			s.Scan(nil)
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
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

	Convey("UnmarshalJSON", t, func() {
		s := Int64{}
		err := json.Unmarshal([]byte("5612"), &s)
		So(err, ShouldBeNil)
		v, err := s.Value()
		So(err, ShouldBeNil)
		So(v, ShouldEqual, int64(5612))
		So(s.IsSet(), ShouldBeTrue)
	})
}

func TestNullInt64(t *testing.T) {
	Convey("Scan", t, func() {
		Convey("Scan should load Int64 and Shadow field", func() {
			ns := &NullInt64{}
			ns.Scan(1234)

			So(ns.Int64, ShouldEqual, 1234)
			So(ns.shadow.Int64, ShouldEqual, 1234)
			So(ns.IsSet(), ShouldBeTrue)
		})
		Convey("secondary Scan should not update shadow", func() {

			ns := &NullInt64{}
			ns.Scan(1234)
			ns.Scan(5678)

			So(ns.Int64, ShouldEqual, 5678)
			So(ns.shadow.Int64, ShouldEqual, 1234)
			So(ns.IsSet(), ShouldBeTrue)
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

		Convey("default to nil", func() {
			ns := &NullInt64{}

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

	Convey("IsSet", t, func() {
		ns := &NullInt64{}
		So(ns.IsSet(), ShouldBeFalse)
		ns.Scan(47473)
		So(ns.IsSet(), ShouldBeTrue)
	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			ns := &NullInt64{}
			value, err := ns.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
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

	Convey("MarshalJSON", t, func() {

		Convey("with value", func() {
			s := NullInt64{}
			s.Scan(1234)
			data, _ := json.Marshal(s)
			So(string(data), ShouldEqual, "1234")
		})

		Convey("with null", func() {
			s := NullInt64{}
			s.Scan(nil)
			data, _ := json.Marshal(s)
			So(string(data), ShouldEqual, "null")
		})

		Convey("default to null when not set", func() {
			s := NullInt64{}
			data, _ := json.Marshal(s)
			So(string(data), ShouldEqual, "null")
		})
	})

	Convey("UnmarshalJSON", t, func() {

		Convey("with value", func() {
			s := NullInt64{}
			err := json.Unmarshal([]byte("5612"), &s)
			So(err, ShouldBeNil)
			v, err := s.Value()
			So(err, ShouldBeNil)
			So(v, ShouldEqual, int64(5612))
			So(s.IsSet(), ShouldBeTrue)
		})

		Convey("with null", func() {
			s := NullInt64{}
			err := json.Unmarshal([]byte("null"), &s)
			So(err, ShouldBeNil)
			_, err = s.Value()
			So(err, ShouldBeNil)
			So(s.Valid, ShouldBeFalse)
			So(s.IsSet(), ShouldBeTrue)
		})
	})
}
