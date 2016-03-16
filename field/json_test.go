package field

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNullJSON(t *testing.T) {
	Convey("NullJSON Scan", t, func() {
		Convey("nil", func() {
			Convey("null as string should have value nil", func() {
				js := &NullJSON{}
				err := js.Scan("null")
				So(err, ShouldBeNil)
				val, _ := js.Value()
				So(val, ShouldBeNil)
			})

			Convey("Scanning nil should return nil", func() {
				js := &NullJSON{}
				err := js.Scan(nil)
				So(err, ShouldBeNil)
				val, _ := js.Value()
				So(val, ShouldBeNil)
			})
		})

		Convey("NullJSON Object", func() {
			js := &NullJSON{}
			err := js.Scan(`{"data1":1, "data2":"hello"}`)
			So(err, ShouldBeNil)
			data1 := js.NullJSON.(map[string]interface{})["data1"]
			data2 := js.NullJSON.(map[string]interface{})["data2"]
			So(data1, ShouldEqual, 1)
			So(data2, ShouldEqual, "hello")
		})

		Convey("NullJSON List", func() {
			js := &NullJSON{}
			err := js.Scan(`[1,2,3]`)
			So(err, ShouldBeNil)
			So(js.NullJSON.([]interface{})[0], ShouldEqual, 1)
			So(js.NullJSON.([]interface{})[1], ShouldEqual, 2)
			So(js.NullJSON.([]interface{})[2], ShouldEqual, 3)
		})

		Convey("NullJSON int", func() {
			js := &NullJSON{}
			err := js.Scan(`154`)
			So(err, ShouldBeNil)
			So(js.NullJSON, ShouldEqual, 154)
		})

		Convey("NullJSON float", func() {
			js := &NullJSON{}
			err := js.Scan(`3.1415926`)
			So(err, ShouldBeNil)
			So(js.NullJSON, ShouldEqual, 3.1415926)
		})

		Convey("NullJSON string", func() {
			js := &NullJSON{}
			err := js.Scan(`"hello world"`)
			So(err, ShouldBeNil)
			So(js.NullJSON, ShouldEqual, "hello world")
		})

		Convey("Fail", func() {
			Convey("string with no quotes", func() {
				js := &NullJSON{}
				err := js.Scan(`hello world`)
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Value", t, func() {
		js := &NullJSON{}
		js.Scan(`{"data":{"attributes":{}}}`)
		val, err := js.Value()
		So(err, ShouldBeNil)
		So(val, ShouldEqual, `{"data":{"attributes":{}}}`)
	})

	Convey("IsDirty", t, func() {
		Convey("changing value should have isdirty true", func() {
			js := &NullJSON{}
			js.Scan(`{"data":{"attributes":{}}}`)
			js.Scan(`{"data":{"attributes":{"name":"helloworld"}}}`)
			So(js.IsDirty(), ShouldBeTrue)
		})

		Convey("two same type values equal", func() {
			js := &NullJSON{}
			js.Scan(`4`)
			js.Scan(`4`)
			So(js.IsDirty(), ShouldBeTrue)
		})

		Convey("one scan should return false", func() {
			js := &NullJSON{}
			js.Scan(`7`)
			So(js.IsDirty(), ShouldBeFalse)
		})
	})
}
