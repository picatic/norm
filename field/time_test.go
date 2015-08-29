package field

import (
	"encoding/json"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05.000"
	timeA      = "2015-01-01 12:12:12.000"
	timeB      = "2014-01-01 12:12:12.000"
)

func TestTime(t *testing.T) {
	Convey("Scan", t, func() {
		Convey("Scan should load Time and Shadow field", func() {
			s := &Time{}
			s.Scan(timeA)

			So(s.Time.Format(timeFormat), ShouldEqual, timeA)
			So(s.shadow.Format(timeFormat), ShouldEqual, timeA)
		})
		Convey("secondary Scan should not update shadow", func() {

			s := &Time{}
			s.Scan(timeA)
			s.Scan(timeB)

			So(s.Time.Format(timeFormat), ShouldEqual, timeB)
			So(s.shadow.Format(timeFormat), ShouldEqual, timeA)
		})
	})

	Convey("Value", t, func() {
		Convey("should return scanned value", func() {
			s := &Time{}
			s.Scan(timeA)

			value, err := s.Value()
			timeValue := value.(time.Time)
			So(timeValue.Format(timeFormat), ShouldEqual, timeA)
			So(err, ShouldBeNil)
		})

		Convey("should return err on scanned nil", func() {
			s := &Time{}
			s.Scan(nil)

			value, err := s.Value()
			So(value, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("IsDirty", t, func() {
		Convey("should be false", func() {
			s := &Time{}
			s.Scan(timeA)

			So(s.IsDirty(), ShouldBeFalse)
		})

		Convey("should be false if nil", func() {
			s := &Time{}
			s.Scan(nil)

			So(s.IsDirty(), ShouldBeFalse)
		})

		Convey("should be true if modified", func() {
			s := &Time{}
			s.Scan(timeA)
			s.Scan(timeB)

			So(s.IsDirty(), ShouldBeTrue)
		})

		Convey("should be false if modified from nil", func() {
			s := &Time{}
			s.Scan(nil)   // does not set to empty string "" as it errors out.
			s.Scan(timeA) // sets to timeB

			So(s.IsDirty(), ShouldBeFalse)
		})

	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			s := &Time{}
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, errors.New("Shadow Wasn't Created"))
		})
		Convey("should return error when only a nil scanned", func() {
			s := &Time{}
			s.Scan(nil)
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, errors.New("Shadow Wasn't Created"))
		})
		Convey("should return scanned Time", func() {
			s := &Time{}
			s.Scan(timeA)
			value, err := s.ShadowValue()
			timeValue := value.(time.Time)
			So(timeValue.Format(timeFormat), ShouldEqual, timeA)
			So(err, ShouldBeNil)
		})
	})

	Convey("MarshalJSON", t, func() {
		s := Time{}
		s.Scan(timeA)
		data, _ := json.Marshal(s)
		So(string(data), ShouldEqual, "\"2015-01-01T12:12:12Z\"")
	})
}

func TestNullTime(t *testing.T) {
	Convey("Scan", t, func() {
		Convey("Scan should load Time and Shadow field", func() {
			ns := &NullTime{}
			ns.Scan(timeA)

			So(ns.Time.Format(timeFormat), ShouldEqual, timeA)
			So(ns.shadow.Time.Format(timeFormat), ShouldEqual, timeA)
		})
		Convey("secondary Scan should not update shadow", func() {

			ns := &NullTime{}
			ns.Scan(timeA)
			ns.Scan(timeB)

			So(ns.Time.Format(timeFormat), ShouldEqual, timeB)
			So(ns.shadow.Time.Format(timeFormat), ShouldEqual, timeA)
		})
	})

	Convey("Value", t, func() {
		Convey("should return scanned value", func() {
			ns := &NullTime{}
			ns.Scan(timeA)

			value, err := ns.Value()
			timeValue := value.(time.Time)
			So(timeValue.Format(timeFormat), ShouldEqual, timeA)
			So(err, ShouldBeNil)
		})

		Convey("should return scanned nil", func() {
			ns := &NullTime{}
			ns.Scan(nil)

			value, err := ns.Value()
			So(value, ShouldBeNil)
			So(err, ShouldBeNil)
		})
	})

	Convey("IsDirty", t, func() {
		Convey("should be false", func() {
			ns := &NullTime{}
			ns.Scan(timeA)

			So(ns.IsDirty(), ShouldBeFalse)
		})
		Convey("should be false if nil", func() {
			ns := &NullTime{}
			ns.Scan(nil)

			So(ns.IsDirty(), ShouldBeFalse)
		})

		Convey("should be true if modified", func() {
			ns := &NullTime{}
			ns.Scan(timeA)
			ns.Scan(timeB)

			So(ns.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified from nil", func() {
			ns := &NullTime{}
			ns.Scan(nil)
			ns.Scan(timeB)

			So(ns.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified to nil", func() {
			ns := &NullTime{}
			ns.Scan(timeA)
			ns.Scan(nil)

			So(ns.IsDirty(), ShouldBeTrue)
		})

	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			ns := &NullTime{}
			value, err := ns.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, errors.New("Shadow Wasn't Created"))
		})
		Convey("should return nil before when nil", func() {
			ns := &NullTime{}
			ns.Scan(nil)
			value, err := ns.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldBeNil)
		})
		Convey("should return scanned string", func() {
			ns := &NullTime{}
			ns.Scan(timeA)
			value, err := ns.ShadowValue()
			timeValue := value.(time.Time)
			So(timeValue.Format(timeFormat), ShouldEqual, timeA)
			So(err, ShouldBeNil)
		})
	})
}
