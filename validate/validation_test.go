package validate

import (
	"testing"

	"github.com/picatic/norm/field"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidate(t *testing.T) {
	Convey("Validation", t, func() {
		Convey("NotNullable", func() {
			Convey("scanning nil should return an error", func() {
				nf := &normFields{}
				nf.NullString.Scan(nil)
				err := NormField("NullString", NotNullable(Always)).Validate(nf)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("NormFields", func() {
			Convey("String", func() {
				str := normFields{}
				str.String.Scan("helloworld")
				err := NormField("String", ValidatorFunc(func(v interface{}) error {
					return nil
				})).Validate(str)
				So(err, ShouldBeNil)
			})
		})

		Convey("IsStringUUID", func() {
			Convey("if string is uuid should not return error", func() {
				str := normFields{}
				str.String.Scan(uuid.NewV4().String())
				err := NormField("String", UUID).Validate(str)
				So(err, ShouldBeNil)
			})

			Convey("if string is not uuid should return error", func() {
				str := normFields{}
				str.String.Scan("steve")
				err := NormField("String", UUID).Validate(str)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("IsStringEmail", func() {
			Convey("non email should be invalid", func() {
				str := normFields{}
				str.String.Scan("steve")
				err := NormField("String", Email).Validate(str)
				So(err, ShouldNotBeNil)
			})

			Convey("email should be valid", func() {
				str := normFields{}
				str.String.Scan("steve@gmail.com")
				err := NormField("String", Email).Validate(str)
				So(err, ShouldBeNil)
			})
		})

		Convey("IsStringInRange", func() {
			Convey("valid", func() {
				str := normFields{}
				str.String.Scan("123456")
				err := NormField("String", Length(All(
					GTE(1),
					LTE(7),
				))).Validate(str)
				So(err, ShouldBeNil)
			})

			Convey("invalid", func() {
				str := normFields{}
				str.String.Scan("12345678")
				err := NormField("String", Length(All(
					GTE(1),
					LTE(7),
				))).Validate(str)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("List", func() {
			InRange0To100 := All(
				GT(0),
				LT(100),
			)

			fl := fieldList{}
			fl.Raws = make([]raw, 3)
			fl.Raws[0].Int64 = 10
			fl.Raws[1].Int64 = 20
			fl.Raws[2].Int64 = 30

			validRawList := Field("Raws", List(Field("Int64", InRange0To100)))
			err := validRawList.Validate(fl)
			So(err, ShouldBeNil)

			Convey("List In List", func() {
				ll := fieldListinList{}
				ll.List = make([]fieldList, 2)

				ll.List[0].Raws = make([]raw, 3)
				ll.List[0].Raws[0].Int64 = 10
				ll.List[0].Raws[1].Int64 = 20
				ll.List[0].Raws[2].Int64 = 30

				ll.List[1].Raws = make([]raw, 3)
				ll.List[1].Raws[0].Int64 = 10
				ll.List[1].Raws[1].Int64 = 20
				ll.List[1].Raws[2].Int64 = 110 //this should cause an error

				err := Field("List", List(validRawList)).Validate(ll)

				So(err, ShouldNotBeNil)
			})
		})
	})
}

type raw struct {
	Int64   int64
	Int32   int32
	Int16   int16
	Int8    int8
	Float64 float64
	Float32 float32
	String  string
}

type normFields struct {
	Int64       field.Int64
	Float64     field.Float64
	String      field.String
	Bool        field.Bool
	NullString  field.NullString
	NullBool    field.NullBool
	NullFloat64 field.NullFloat64
	NullInt64   field.NullInt64
}

type fieldList struct {
	Raws       []raw
	NormFields []normFields
}

type fieldListinList struct {
	List []fieldList
}
