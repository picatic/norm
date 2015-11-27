package field

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

const (
	timeTimeA = "12:12:13"
	timeTimeB = "12:12:12"
)

func TestTimeTime(t *testing.T) {
	var (
		timeStructA time.Time
	)
	timeStructA, _ = time.Parse(timeTimeFormat, timeTimeA)
	Convey("Scan", t, func() {
		Convey("Scan should load Time and Shadow field", func() {
			s := &TimeTime{}
			s.Scan(timeTimeA)

			So(s.Time.Format(timeTimeFormat), ShouldEqual, timeTimeA)
			So(s.shadow.Format(timeTimeFormat), ShouldEqual, timeTimeA)
			So(s.IsSet(), ShouldBeTrue)
		})
		Convey("secondary Scan should not update shadow", func() {

			s := &TimeTime{}
			s.Scan(timeTimeA)
			s.Scan(timeTimeB)

			So(s.Time.Format(timeTimeFormat), ShouldEqual, timeTimeB)
			So(s.shadow.Format(timeTimeFormat), ShouldEqual, timeTimeA)
			So(s.IsSet(), ShouldBeTrue)
		})

		Convey("can set with time.Time", func() {
			s := &TimeTime{}
			err := s.Scan(timeStructA)

			So(err, ShouldBeNil)

			So(s.Time.Format(timeTimeFormat), ShouldEqual, timeTimeA)
			So(s.IsSet(), ShouldBeTrue)
		})

		Convey("Zero time is an error", func() {
			s := &TimeTime{}
			t := time.Time{}
			err := s.Scan(t)

			So(err, ShouldNotBeNil)

			So(s.IsSet(), ShouldBeFalse)
		})

		Convey("can set with bytes", func() {
			s := &TimeTime{}
			err := s.Scan([]byte(timeTimeA))

			So(err, ShouldBeNil)

			So(s.Time.Format(timeTimeFormat), ShouldEqual, timeTimeA)
			So(s.IsSet(), ShouldBeTrue)
		})
	})

	Convey("Value", t, func() {
		Convey("should return scanned value", func() {
			s := &TimeTime{}
			s.Scan(timeTimeA)

			value, err := s.Value()
			timeValue := value.(time.Time)
			So(timeValue.Format(timeTimeFormat), ShouldEqual, timeTimeA)
			So(err, ShouldBeNil)
		})

		Convey("should return Zero time on scanned nil", func() {
			s := &TimeTime{}
			s.Scan(nil)

			value, err := s.Value()
			So(value.(time.Time).IsZero(), ShouldBeTrue)
			So(err, ShouldBeNil)
		})

		Convey("should return Zero time when not set", func() {
			s := &TimeTime{}

			value, err := s.Value()
			So(value.(time.Time).IsZero(), ShouldBeTrue)
			So(err, ShouldBeNil)
		})
	})

	Convey("IsDirty", t, func() {
		Convey("should be false", func() {
			s := &TimeTime{}
			s.Scan(timeTimeA)

			So(s.IsDirty(), ShouldBeFalse)
		})

		Convey("should be false if nil", func() {
			s := &TimeTime{}
			s.Scan(nil)

			So(s.IsDirty(), ShouldBeFalse)
		})

		Convey("should be true if modified", func() {
			s := &TimeTime{}
			s.Scan(timeTimeA)
			s.Scan(timeTimeB)

			So(s.IsDirty(), ShouldBeTrue)
		})

		Convey("should be false if modified from nil", func() {
			s := &TimeTime{}
			s.Scan(nil)       // does not set to empty string "" as it errors out.
			s.Scan(timeTimeA) // sets to timeTimeB

			So(s.IsDirty(), ShouldBeFalse)
		})
	})

	Convey("IsSet", t, func() {
		s := TimeTime{}
		So(s.IsSet(), ShouldBeFalse)
		s.Scan(timeTimeA)
		So(s.IsSet(), ShouldBeTrue)
	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			s := &TimeTime{}
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
		})
		Convey("should return error when only a nil scanned", func() {
			s := &TimeTime{}
			s.Scan(nil)
			value, err := s.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
		})
		Convey("should return scanned Time", func() {
			s := &TimeTime{}
			s.Scan(timeTimeA)
			value, err := s.ShadowValue()
			timeValue := value.(time.Time)
			So(timeValue.Format(timeTimeFormat), ShouldEqual, timeTimeA)
			So(err, ShouldBeNil)
		})
	})

	Convey("MarshalJSON", t, func() {
		s := TimeTime{}
		s.Scan(timeTimeA)
		data, _ := json.Marshal(s)
		So(string(data), ShouldEqual, "\"12:12:13\"")
	})

	Convey("UnmarshalJSON", t, func() {
		s := TimeTime{}
		err := json.Unmarshal([]byte("\"12:12:12\""), &s)
		So(err, ShouldBeNil)
		v, err := s.Value()
		So(err, ShouldBeNil)
		So(v.(time.Time).String(), ShouldEqual, "0000-01-01 12:12:12 +0000 UTC")
		So(s.IsSet(), ShouldBeTrue)
	})
}

func TestNullTimeTime(t *testing.T) {
	Convey("Scan", t, func() {
		Convey("Scan should load Time and Shadow field", func() {
			ns := NullTimeTime{}
			err := ns.Scan(timeTimeA)
			So(err, ShouldBeNil)
			So(ns.Time.Format(timeTimeFormat), ShouldEqual, timeTimeA)
			So(ns.shadow.Time.Format(timeTimeFormat), ShouldEqual, timeTimeA)
			So(ns.IsSet(), ShouldBeTrue)
		})

		Convey("zero time is 12:00:00AM", func() {
			ns := NullTimeTime{}
			t := time.Time{}
			err := ns.Scan(t)
			So(err, ShouldBeNil)
			So(ns.Valid, ShouldBeTrue)
			So(ns.IsSet(), ShouldBeTrue)
		})
		Convey("secondary Scan should not update shadow", func() {

			ns := &NullTimeTime{}
			ns.Scan(timeTimeA)
			ns.Scan(timeTimeB)

			So(ns.Time.Format(timeTimeFormat), ShouldEqual, timeTimeB)
			So(ns.shadow.Time.Format(timeTimeFormat), ShouldEqual, timeTimeA)
			So(ns.IsSet(), ShouldBeTrue)
		})
	})

	Convey("Value", t, func() {
		Convey("should return scanned value", func() {
			ns := &NullTimeTime{}
			ns.Scan(timeTimeA)

			value, err := ns.Value()
			So(err, ShouldBeNil)
			timeValue := value.(time.Time)
			So(timeValue.Format(timeTimeFormat), ShouldEqual, timeTimeA)
			So(err, ShouldBeNil)
		})

		Convey("should return scanned nil", func() {
			ns := &NullTimeTime{}
			ns.Scan(nil)

			value, err := ns.Value()
			So(err, ShouldBeNil)
			So(value, ShouldBeNil)
		})
	})

	Convey("IsDirty", t, func() {
		Convey("should be false", func() {
			ns := &NullTimeTime{}
			err := ns.Scan(timeTimeA)
			So(err, ShouldBeNil)

			So(ns.IsDirty(), ShouldBeFalse)
		})

		Convey("should be false if nil", func() {
			ns := &NullTimeTime{}
			ns.Scan(nil)

			So(ns.IsDirty(), ShouldBeFalse)
		})

		Convey("should be true if modified", func() {
			ns := &NullTimeTime{}
			ns.Scan(timeTimeA)
			ns.Scan(timeTimeB)

			So(ns.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified from nil", func() {
			ns := &NullTimeTime{}
			err := ns.Scan(nil)
			So(err, ShouldBeNil)
			ns.Scan(timeTimeB)

			So(ns.IsDirty(), ShouldBeTrue)
		})
		Convey("should be true if modified to nil", func() {
			ns := &NullTimeTime{}
			ns.Scan(timeTimeA)
			ns.Scan(nil)

			So(ns.IsDirty(), ShouldBeTrue)
		})

	})

	Convey("IsSet", t, func() {
		s := NullTimeTime{}
		So(s.IsSet(), ShouldBeFalse)
		s.Scan(timeTimeA)
		So(s.IsSet(), ShouldBeTrue)
	})

	Convey("ShadowValue", t, func() {
		Convey("should return err before first scan", func() {
			ns := &NullTimeTime{}
			value, err := ns.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldResemble, ErrorUnintializedShadow)
		})
		Convey("should return nil before when nil", func() {
			ns := &NullTimeTime{}
			ns.Scan(nil)
			value, err := ns.ShadowValue()

			So(value, ShouldBeNil)
			So(err, ShouldBeNil)
		})
		Convey("should return scanned string", func() {
			ns := &NullTimeTime{}
			ns.Scan(timeTimeA)
			value, err := ns.ShadowValue()
			timeValue := value.(time.Time)
			So(timeValue.Format(timeTimeFormat), ShouldEqual, timeTimeA)
			So(err, ShouldBeNil)
		})
	})

	Convey("MarshalJSON", t, func() {

		Convey("valid date", func() {
			s := NullTimeTime{}
			s.Scan(timeTimeA)
			data, _ := json.Marshal(s)
			So(string(data), ShouldEqual, "\"12:12:13\"")
		})

		Convey("null date", func() {
			s := NullTimeTime{}
			s.Scan(nil)
			data, _ := json.Marshal(s)
			So(string(data), ShouldEqual, "null")
		})

	})

	Convey("UnmarshalJSON", t, func() {
		Convey("Valid TimeTime", func() {
			s := NullTimeTime{}
			err := json.Unmarshal([]byte("\"12:12:12\""), &s)
			So(err, ShouldBeNil)
			v, err := s.Value()
			So(err, ShouldBeNil)
			So(v.(time.Time).String(), ShouldEqual, "0000-01-01 12:12:12 +0000 UTC")
			So(s.IsSet(), ShouldBeTrue)
		})

		Convey("Invalid Date format", func() {
			s := NullTimeTime{}
			err := json.Unmarshal([]byte("\"2015-01-01 12:12:12\""), &s)
			So(err, ShouldNotBeNil)
			So(s.IsSet(), ShouldBeFalse)
		})

		Convey("Invalid empty string", func() {
			s := NullTimeTime{}
			err := json.Unmarshal([]byte("\"\""), &s)
			So(err, ShouldNotBeNil)
			So(s.IsSet(), ShouldBeFalse)
		})

		Convey("00:00:00 is valid", func() {
			s := NullTimeTime{}
			err := json.Unmarshal([]byte("\"00:00:00\""), &s)
			So(err, ShouldBeNil)
			So(s.Valid, ShouldBeTrue)
			So(s.IsSet(), ShouldBeTrue)
		})

		Convey("null date", func() {
			s := NullTimeTime{}
			err := json.Unmarshal([]byte("null"), &s)
			So(err, ShouldBeNil)
			v, err := s.Value()
			So(err, ShouldBeNil)
			So(v, ShouldBeNil)
			So(s.IsSet(), ShouldBeTrue)
		})

	})
}
