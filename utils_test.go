package norm

import (
	"github.com/picatic/go-api/norm/field"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUtils(t *testing.T) {
	Convey("Utils", t, func() {
		model := &MockModel{}
		model.Id.Scan("1")
		model.FirstName.Scan("A First Name")

		Convey("escapeFields", func() {
			fns := field.FieldNames{"Id", "FirstName"}
			So(escapeFields(fns), ShouldResemble, []string{"`id`", "`first_name`"})
		})

		Convey("defaultFieldEscaped", nil)

		Convey("defaultUpdate", func() {
			fns := field.FieldNames{"FirstName"}
			mapSet := defaultUpdate(model, fns)
			ns := field.NullString{}
			ns.Scan("A First Name")
			So(mapSet, ShouldResemble, map[string]interface{}{"first_name": ns})
		})
	})
}
