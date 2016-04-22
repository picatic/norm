package field

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

const (
	timeDateA = "2015-01-01"
	timeDateB = "2014-01-01"
)

func TestTimeDate(t *testing.T) {
	var (
		timeStructA time.Time
	)
	timeStructA, _ = time.Parse(timeDateFormat, timeDateA)
	Convey("Scan", t, func() {
		Convey("Scan should load Time and Shadow field", func() {
			s := &TimeDate{}
			s.Scan(timeDateA)

			So(s.Time.Format(timeDateFormat), ShouldEqual, timeDateA)
			So(s.shadow.Format(timeDateFormat), ShouldEqual, timeDateA)
			So(s.IsSet(), ShouldBeTrue)
		})
		Convey("secondary Scan should not update shadow", func() {

			s := &TimeDate{}
			s.Scan(timeDateA)
			s.Scan(timeDateB)

			So(s.Time.Format(timeDateFormat), ShouldEqual, timeDateB)
			So(s.shadow.Format(timeDateFormat), ShouldEqual, timeDateA)
			So(s.IsSet(), ShouldBeTrue)
		})

		Convey("can set with time.Time", func() {
			s := &TimeDate{}
			err := s.Scan(timeStructA)

			So(err, ShouldBeNil)

			So(s.Time.Format(timeDateFormat), ShouldEqual, timeDateA)
			So(s.IsSet(), ShouldBeTrue)
		})

		Convey("Zero time is an error", func() {
			s := &TimeDate{}
			t := time.Time{}
			err := s.Scan(t)

			So(err, ShouldNotBeNil)

			So(s.IsSet(), ShouldBeFalse)
		})

		Convey("can set with bytes", func() {
			s := &TimeDate{}
			err := s.Scan([]byte(timeDateA))

			So(err, ShouldBeNil)

			So(s.Time.Format(timeDateFormat), ShouldEqual, timeDateA)
			So(s.IsSet(), ShouldBeTrue)
		})
	})

	Convey("Value", t, func() {
		Convey("should return scanned value", func() {
			s := &TimeDate{}
			s.Scan(timeDateA)

			value, err := s.Value()
			timeValue := value.(time.Time)
			So(timeValue.Format(timeDateFormat), ShouldEqual, timeDateA)
			So(err, ShouldBeNil)
		})

		Convey("should return Zero time on scanned nil", func() {
			s := &TimeDate{}
			s.Scan(nil)

			value, err := s.Value()
			So(value.(time.Time).IsZero(), ShouldBeTrue)
			So(err, ShouldBeNil)
		})

		Convey("should return Zero time when not set", func() {
			s := &TimeDate{}

			value, err := s.Value()
			So(value.(time.Time).IsZero(), ShouldBeTrue)
			So(err, ShouldBeNil)
		})
	})

	Convey("IsDirty", t, func() {
		Convey("should be false", func() {
			s := &TimeDate{}
			s.Scan(timeDateA)

			So(s.IsDirty(), ShouldBeFalse)
		})

		Convey("should be false if nil", func() {
			s := &TimeDate{}
			s.Scan(nil)

			So(s.IsDirty(), ShouldBeFalse)
		})

		Convey("should be true if modified", func() {
			s := &TimeDate{}
			s.Scan(timeDateA)
			s.Scan(timeDateB)

			So(s.IsDirty(), ShouldBeTrue)
		})

		Convey("should be false if modified from nil", func() {
			s := &TimeDate{}
			s.Scan(nil)       // does not set to empty string "" as it errors out.
			s.Scan(timeDateA) // sets to timeDateB

			So(s.IsDirty(), ShouldBeFalse)
		})
	})

	Convey("IsSet", t, func() {
		s := TimeDate{}
		So(s.IsSet(), ShouldBeFalse)
		s.Scan(timeDateA)
		So(s.IsSet(), ShouldBeTrue)
	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			s := &TimeDate{}
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
		})
		Convey("should return error when only a nil scanned", func() {
			s := &TimeDate{}
			s.Scan(nil)
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
		})
		Convey("should return scanned Time", func() {
			s := &TimeDate{}
			s.Scan(timeDateA)
			value, err := s.ShadowValue()
			timeValue := value.(time.Time)
			So(timeValue.Format(timeDateFormat), ShouldEqual, timeDateA)
			So(err, ShouldBeNil)
		})
	})

	Convey("MarshalJSON", t, func() {
		s := TimeDate{}
		s.Scan(timeDateA)
		data, _ := json.Marshal(s)
		So(string(data), ShouldEqual, "\"2015-01-01\"")
	})

	Convey("UnmarshalJSON", t, func() {
		s := TimeDate{}
		err := json.Unmarshal([]byte("\"2015-01-01\""), &s)
		So(err, ShouldBeNil)
		v, err := s.Value()
		So(err, ShouldBeNil)
		So(v.(time.Time).String(), ShouldEqual, "2015-01-01 00:00:00 +0000 UTC")
		So(s.IsSet(), ShouldBeTrue)
	})
}

func TestNullTimeDate(t *testing.T) {
	Convey("Unscanned", t, func() {
		Convey("Value should be Null", func() {
			ns := NullTimeDate{}
			t, err := ns.Value()
			So(err, ShouldBeNil)
			So(t, ShouldBeNil)
		})

		Convey("IsDirty should be false", func() {
			ns := NullTimeDate{}
			So(ns.IsDirty(), ShouldBeFalse)
		})

		Convey("Marshal should provide json null", func() {
			ns := NullTimeDate{}
			v, err := ns.MarshalJSON()
			So(err, ShouldBeNil)
			So(string(v), ShouldEqual, "null")
		})
	})

	Convey("Scan", t, func() {
		Convey("Scan should load Time and Shadow field", func() {
			ns := NullTimeDate{}
			err := ns.Scan(timeDateA)
			So(err, ShouldBeNil)
			So(ns.Time.Format(timeDateFormat), ShouldEqual, timeDateA)
			So(ns.shadow.Time.Format(timeDateFormat), ShouldEqual, timeDateA)
			So(ns.IsSet(), ShouldBeTrue)
		})

		Convey("zero time is nil", func() {
			ns := NullTimeDate{}
			t := time.Time{}
			err := ns.Scan(t)
			So(err, ShouldBeNil)
			So(ns.Valid, ShouldBeFalse)
			So(ns.IsSet(), ShouldBeTrue)
		})
		Convey("secondary Scan should not update shadow", func() {

			ns := &NullTimeDate{}
			ns.Scan(timeDateA)
			ns.Scan(timeDateB)

			So(ns.Time.Format(timeDateFormat), ShouldEqual, timeDateB)
			So(ns.shadow.Time.Format(timeDateFormat), ShouldEqual, timeDateA)
			So(ns.IsSet(), ShouldBeTrue)
		})
	})

	Convey("Value", t, func() {
		Convey("should return scanned value", func() {
			ns := &NullTimeDate{}
			ns.Scan(timeDateA)

			value, err := ns.Value()
			So(err, ShouldBeNil)
			timeValue := value.(time.Time)
			So(timeValue.Format(timeDateFormat), ShouldEqual, timeDateA)
			So(err, ShouldBeNil)
		})

		Convey("should return scanned nil", func() {
			ns := &NullTimeDate{}
			ns.Scan(nil)

			value, err := ns.Value()
			So(err, ShouldBeNil)
			So(value, ShouldBeNil)
		})
	})

	Convey("IsDirty", t, func() {
		Convey("should be false", func() {
			ns := &NullTimeDate{}
			err := ns.Scan(timeDateA)
			So(err, ShouldBeNil)

			So(ns.IsDirty(), ShouldBeFalse)
		})

		Convey("should be false if nil", func() {
			ns := &NullTimeDate{}
			ns.Scan(nil)

			So(ns.IsDirty(), ShouldBeFalse)
		})

		Convey("should be true if modified", func() {
			ns := &NullTimeDate{}
			ns.Scan(timeDateA)
			ns.Scan(timeDateB)

			So(ns.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified from nil", func() {
			ns := &NullTimeDate{}
			err := ns.Scan(nil)
			So(err, ShouldBeNil)
			ns.Scan(timeDateB)

			So(ns.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified to nil", func() {
			ns := &NullTimeDate{}
			ns.Scan(timeDateA)
			ns.Scan(nil)

			So(ns.IsDirty(), ShouldBeTrue)
		})

	})

	Convey("IsSet", t, func() {
		s := NullTimeDate{}
		So(s.IsSet(), ShouldBeFalse)
		s.Scan(timeDateA)
		So(s.IsSet(), ShouldBeTrue)
	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			ns := &NullTimeDate{}
			value, err := ns.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
		})
		Convey("should return nil before when nil", func() {
			ns := &NullTimeDate{}
			ns.Scan(nil)
			value, err := ns.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldBeNil)
		})
		Convey("should return scanned string", func() {
			ns := &NullTimeDate{}
			ns.Scan(timeDateA)
			value, err := ns.ShadowValue()
			timeValue := value.(time.Time)
			So(timeValue.Format(timeDateFormat), ShouldEqual, timeDateA)
			So(err, ShouldBeNil)
		})
	})

	Convey("MarshalJSON", t, func() {

		Convey("valid date", func() {
			s := NullTimeDate{}
			s.Scan(timeDateA)
			data, _ := json.Marshal(s)
			So(string(data), ShouldEqual, "\"2015-01-01\"")
		})

		Convey("null date", func() {
			s := NullTimeDate{}
			s.Scan(nil)
			data, _ := json.Marshal(s)
			So(string(data), ShouldEqual, "null")
		})

	})

	Convey("UnmarshalJSON", t, func() {
		Convey("Valid Date", func() {
			s := NullTimeDate{}
			err := json.Unmarshal([]byte("\"2015-01-01\""), &s)
			So(err, ShouldBeNil)
			v, err := s.Value()
			So(err, ShouldBeNil)
			So(v.(time.Time).String(), ShouldEqual, "2015-01-01 00:00:00 +0000 UTC")
			So(s.IsSet(), ShouldBeTrue)
		})

		Convey("Invalid Date format", func() {
			s := NullTimeDate{}
			err := json.Unmarshal([]byte("\"2015-01-01-\""), &s)
			So(err, ShouldNotBeNil)
			So(s.IsSet(), ShouldBeFalse)
		})

		Convey("Invalid empty string", func() {
			s := NullTimeDate{}
			err := json.Unmarshal([]byte("\"\""), &s)
			So(err, ShouldNotBeNil)
			So(s.IsSet(), ShouldBeFalse)
		})

		Convey("Invalid date range. as zero time?", func() {
			s := NullTimeDate{}
			err := json.Unmarshal([]byte("\"0001-01-01\""), &s)
			So(err, ShouldBeNil)
			So(s.Valid, ShouldBeFalse)
			So(s.IsSet(), ShouldBeTrue)
		})

		Convey("null date", func() {
			s := NullTimeDate{}
			err := json.Unmarshal([]byte("null"), &s)
			So(err, ShouldBeNil)
			v, err := s.Value()
			So(err, ShouldBeNil)
			So(v, ShouldBeNil)
			So(s.IsSet(), ShouldBeTrue)
		})

	})
}
