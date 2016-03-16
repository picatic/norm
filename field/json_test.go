package field

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNullJson(t *testing.T) {
	Convey("NullJson Scan", t, func() {
		Convey("nil", func() {
			Convey("null as string should have value nil", func() {
				js := &NullJson{}
				err := js.Scan("null")
				So(err, ShouldBeNil)
				val, _ := js.Value()
				So(val, ShouldBeNil)
			})

			Convey("Scanning nil should return nil", func() {
				js := &NullJson{}
				err := js.Scan(nil)
				So(err, ShouldBeNil)
				val, _ := js.Value()
				So(val, ShouldBeNil)
			})
		})

		Convey("NullJson Object", func() {
			js := &NullJson{}
			err := js.Scan(`{"data1":1, "data2":"hello"}`)
			So(err, ShouldBeNil)
			data1 := js.NullJson.(map[string]interface{})["data1"]
			data2 := js.NullJson.(map[string]interface{})["data2"]
			So(data1, ShouldEqual, 1)
			So(data2, ShouldEqual, "hello")
		})

		Convey("NullJson List", func() {
			js := &NullJson{}
			err := js.Scan(`[1,2,3]`)
			So(err, ShouldBeNil)
			So(js.NullJson.([]interface{})[0], ShouldEqual, 1)
			So(js.NullJson.([]interface{})[1], ShouldEqual, 2)
			So(js.NullJson.([]interface{})[2], ShouldEqual, 3)
		})

		Convey("NullJson int", func() {
			js := &NullJson{}
			err := js.Scan(`154`)
			So(err, ShouldBeNil)
			So(js.NullJson, ShouldEqual, 154)
		})

		Convey("NullJson float", func() {
			js := &NullJson{}
			err := js.Scan(`3.1415926`)
			So(err, ShouldBeNil)
			So(js.NullJson, ShouldEqual, 3.1415926)
		})

		Convey("NullJson string", func() {
			js := &NullJson{}
			err := js.Scan(`"hello world"`)
			So(err, ShouldBeNil)
			So(js.NullJson, ShouldEqual, "hello world")
		})

		Convey("Fail", func() {
			Convey("string with no quotes", func() {
				js := &NullJson{}
				err := js.Scan(`hello world`)
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Value", t, func() {
		js := &NullJson{}
		js.Scan(`{"data":{"attributes":{}}}`)
		val, err := js.Value()
		So(err, ShouldBeNil)
		So(val, ShouldEqual, `{"data":{"attributes":{}}}`)
	})

	Convey("IsDirty", t, func() {
		Convey("changing value should have isdirty true", func() {
			js := &NullJson{}
			js.Scan(`{"data":{"attributes":{}}}`)
			js.Scan(`{"data":{"attributes":{"name":"helloworld"}}}`)
			So(js.IsDirty(), ShouldBeTrue)
		})

		Convey("two same type values equal", func() {
			js := &NullJson{}
			js.Scan(`4`)
			js.Scan(`4`)
			So(js.IsDirty(), ShouldBeTrue)
		})

		Convey("one scan should return false", func() {
			js := &NullJson{}
			js.Scan(`7`)
			So(js.IsDirty(), ShouldBeFalse)
		})
	})
}
