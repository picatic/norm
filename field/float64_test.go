package field

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFloat64(t *testing.T) {
	Convey("Scan", t, func() {
		Convey("Scan should load Float64 and Shadow field", func() {
			s := &Float64{}
			s.Scan(1234)

			So(s.Float64, ShouldEqual, 1234)
			So(s.shadow, ShouldEqual, 1234)
			So(s.IsSet(), ShouldBeTrue)
		})
		Convey("secondary Scan should not update shadow", func() {

			s := &Float64{}
			s.Scan(12.34)
			s.Scan(56.78)

			So(s.Float64, ShouldEqual, 56.78)
			So(s.shadow, ShouldEqual, 12.34)
			So(s.IsSet(), ShouldBeTrue)
		})

		Convey("empty string", func() {
			s := &NullFloat64{}

			err := s.Scan("")
			So(err, ShouldNotBeNil)

			v, err := s.Value()
			So(err, ShouldBeNil)
			So(v, ShouldBeNil)
		})
	})

	Convey("Value", t, func() {
		Convey("should return scanned value", func() {
			s := &Float64{}
			s.Scan(12.34)

			value, err := s.Value()

			So(value, ShouldEqual, 12.34)
			So(err, ShouldBeNil)
		})

		Convey("should return default Float64 on scanned nil", func() {
			s := &Float64{}
			s.Scan(nil)

			value, err := s.Value()
			So(value, ShouldEqual, float64(0))
			So(err, ShouldBeNil)
		})
	})

	Convey("IsDirty", t, func() {
		Convey("should be false", func() {
			s := &Float64{}
			s.Scan(12.34)

			So(s.IsDirty(), ShouldBeFalse)
		})
		Convey("should be false if nil", func() {
			s := &Float64{}
			s.Scan(nil)

			So(s.IsDirty(), ShouldBeFalse)
		})

		Convey("should be true if modified", func() {
			s := &Float64{}
			s.Scan(12.34)
			s.Scan(56.78)

			So(s.IsDirty(), ShouldBeTrue)
		})
	})

	Convey("IsSet", t, func() {
		s := &Float64{}
		So(s.IsSet(), ShouldBeFalse)
		s.Scan(1.23)
		So(s.IsSet(), ShouldBeTrue)
	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			s := &Float64{}
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
		})
		Convey("should return error when only a nil scanned", func() {
			s := &Float64{}
			s.Scan(nil)
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
		})
		Convey("should return scanned Float64", func() {
			s := &Float64{}
			s.Scan(12.34)
			value, err := s.ShadowValue()

			So(value, ShouldEqual, 12.34)
			So(err, ShouldBeNil)
		})
	})

	Convey("MarshalJSON", t, func() {
		s := Float64{}
		s.Scan(12.34)
		data, _ := json.Marshal(s)
		So(string(data), ShouldEqual, "12.34")
	})

	Convey("UnmarshalJSON", t, func() {
		s := Float64{}
		err := json.Unmarshal([]byte("56.12"), &s)
		So(err, ShouldBeNil)
		v, err := s.Value()
		So(err, ShouldBeNil)
		So(v, ShouldEqual, float64(56.12))
		So(s.IsSet(), ShouldBeTrue)
	})
}

func TestNullFloat64(t *testing.T) {
	Convey("Scan", t, func() {
		Convey("Scan should load Float64 and Shadow field", func() {
			ns := &NullFloat64{}
			ns.Scan(12.34)

			So(ns.Float64, ShouldEqual, 12.34)
			So(ns.shadow.Float64, ShouldEqual, 12.34)
		})
		Convey("secondary Scan should not update shadow", func() {

			ns := &NullFloat64{}
			ns.Scan(12.34)
			ns.Scan(56.78)

			So(ns.Float64, ShouldEqual, 56.78)
			So(ns.shadow.Float64, ShouldEqual, 12.34)
		})
	})

	Convey("Value", t, func() {
		Convey("should return scanned value", func() {
			ns := &NullFloat64{}
			ns.Scan(12.34)

			value, err := ns.Value()

			So(value, ShouldEqual, 12.34)
			So(err, ShouldBeNil)
		})

		Convey("should return scanned nil", func() {
			ns := &NullFloat64{}
			ns.Scan(nil)

			value, err := ns.Value()
			So(value, ShouldBeNil)
			So(err, ShouldBeNil)
		})
	})

	Convey("IsDirty", t, func() {
		Convey("should be false", func() {
			ns := &NullFloat64{}
			ns.Scan(12.34)

			So(ns.IsDirty(), ShouldBeFalse)
		})
		Convey("should be false if nil", func() {
			ns := &NullFloat64{}
			ns.Scan(nil)

			So(ns.IsDirty(), ShouldBeFalse)
		})

		Convey("should be true if modified", func() {
			ns := &NullFloat64{}
			ns.Scan(12.34)
			ns.Scan(56.78)

			So(ns.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified from nil", func() {
			ns := &NullFloat64{}
			ns.Scan(nil)
			ns.Scan(56.78)

			So(ns.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified to nil", func() {
			ns := &NullFloat64{}
			ns.Scan(12.34)
			ns.Scan(nil)

			So(ns.IsDirty(), ShouldBeTrue)
		})

	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			ns := &NullFloat64{}
			value, err := ns.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
		})
		Convey("should return nil before when nil", func() {
			ns := &NullFloat64{}
			ns.Scan(nil)
			value, err := ns.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldBeNil)
		})
		Convey("should return scanned Float64", func() {
			ns := &NullFloat64{}
			ns.Scan(12.34)
			value, err := ns.ShadowValue()

			So(value, ShouldEqual, 12.34)
			So(err, ShouldBeNil)
		})
	})

	Convey("MarshalJSON", t, func() {

		Convey("with value", func() {
			s := NullFloat64{}
			s.Scan(12.34)
			data, _ := json.Marshal(s)
			So(string(data), ShouldEqual, "12.34")
		})

		Convey("with null", func() {
			s := NullFloat64{}
			data, _ := json.Marshal(s)
			So(string(data), ShouldEqual, "null")
		})
	})

	Convey("UnmarshalJSON", t, func() {

		Convey("with value", func() {
			s := NullFloat64{}
			err := json.Unmarshal([]byte("56.12"), &s)
			So(err, ShouldBeNil)
			v, err := s.Value()
			So(err, ShouldBeNil)
			So(v, ShouldEqual, float64(56.12))
			So(s.IsSet(), ShouldBeTrue)
		})

		Convey("with null", func() {
			s := NullFloat64{}
			err := json.Unmarshal([]byte("null"), &s)
			So(err, ShouldBeNil)
			_, err = s.Value()
			So(err, ShouldBeNil)
			So(s.Valid, ShouldBeFalse)
			So(s.IsSet(), ShouldBeTrue)
		})
	})
}
