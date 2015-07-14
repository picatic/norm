package field

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)



func TestOnceDone(t *testing.T) {
	Convey("OnceDone Do", t, func() {

		counter := 0
		once1 := new(ShadowInit)

		Convey("Once Should be called first time", func() {
			// increment counter on first call
			once1.DoInit(func() { counter++ })
			So(counter, ShouldEqual, 1)

		})

		Convey("Incrementer should not be incremented on subsequent calls", func() {
			// do not increment counter on next call
			once1.DoInit(func() { counter++ })
			So(counter, ShouldEqual, 1)

		})

	})
	Convey("OnceDone Done", t, func() {

		counter := 0
		once1 := new(ShadowInit)

		Convey("Done Should return false if Do has not been called", func() {
			// increment counter on first call

			So(once1.InitDone(), ShouldBeFalse)

		})

		Convey("Done should return true if Do has been called", func() {
			// do not increment counter on next call
			once1.DoInit(func() { counter++ })
			So(once1.InitDone(), ShouldBeTrue)

		})

	})
}