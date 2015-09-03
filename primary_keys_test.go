package norm

import (
	"github.com/picatic/norm/field"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPrimaryKeyer(t *testing.T) {
	Convey("PrimaryKeyer", t, func() {
		model := &MockModel{}

		Convey("SinglePrimaryKey", func() {
			pk := NewSinglePrimaryKey(field.Name("Id"))
			So(pk.Fields(), ShouldResemble, field.Names{field.Name("Id")})
			fields, err := model.PrimaryKey().Generator(model)
			So(len(fields), ShouldEqual, 0)
			So(err, ShouldBeNil)
			v, _ := model.Id.Value()
			So(v, ShouldBeNil)
		})

		Convey("MultiplePrimaryKey", func() {
			pk := NewMultiplePrimaryKey(field.Names{"Id", "Org"})
			So(pk.Fields(), ShouldResemble, field.Names{"Id", "Org"})
			fields, err := model.PrimaryKey().Generator(model)
			So(len(fields), ShouldEqual, 0)
			So(err, ShouldBeNil)
			v, _ := model.Id.Value()
			So(v, ShouldBeNil)
			v, _ = model.Org.Value()
			So(v, ShouldBeNil)
		})

		Convey("CustomPrimaryKey", func() {
			pk := NewCustomPrimaryKey(field.Names{"Id", "Org"}, func(pk PrimaryKeyer, model Model) (field.Names, error) {
				f, _ := ModelGetField(model, "Id")
				f.Scan("abc-123")
				return field.Names{"Id"}, nil
			})
			fields, err := pk.Generator(model)
			So(len(fields), ShouldEqual, 1)
			So(err, ShouldBeNil)
			v, _ := model.Id.Value()
			So(v, ShouldEqual, "abc-123")
		})

	})
}
