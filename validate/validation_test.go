package validate

import (
	// "fmt"
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
				err := NormField("NullString", true, Always).Validate(nf)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("NormFields", func() {
			Convey("String", func() {
				str := normFields{}
				str.String.Scan("helloworld")
				err := Field("String", Always).Validate(str)
				So(err, ShouldBeNil)
			})

			Convey("Zero Values", func() {
				fields := &normFields{}

				Convey("Int64", func() {
					validator := NormField("NullInt64", false, Never)

					fields.NullInt64.Scan("")

					err := validator.Validate(fields)
					So(err, ShouldBeNil)
				})

				Convey("Float64", func() {
					validator := NormField("NullFloat64", false, Never)

					fields.NullFloat64.Scan("")

					err := validator.Validate(fields)
					So(err, ShouldBeNil)
				})

				Convey("String", func() {
					validator := NormField("NullString", false, Never)

					fields.NullString.Scan("")

					err := validator.Validate(fields)
					So(err, ShouldBeNil)
				})

				Convey("Bool", func() {
					validator := NormField("NullBool", false, Never)

					fields.NullBool.Scan("")

					err := validator.Validate(fields)
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("IsStringUUID", func() {
			Convey("if string is uuid should not return error", func() {
				str := normFields{}
				str.String.Scan(uuid.NewV4().String())
				err := NormField("String", true, UUID).Validate(str)
				So(err, ShouldBeNil)
			})

			Convey("if string is not uuid should return error", func() {
				str := normFields{}
				str.String.Scan("steve")
				err := NormField("String", true, UUID).Validate(str)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("IsStringEmail", func() {
			Convey("non email should be invalid", func() {
				str := normFields{}
				str.String.Scan("steve")
				err := NormField("String", true, Email).Validate(str)
				So(err, ShouldNotBeNil)
			})

			Convey("email should be valid", func() {
				str := normFields{}
				str.String.Scan("steve@gmail.com")
				err := NormField("String", true, Email).Validate(str)
				So(err, ShouldBeNil)
			})
		})

		Convey("IsStringInRange", func() {
			Convey("valid", func() {
				str := normFields{}
				str.String.Scan("123456")
				err := NormField("String", true, Length(All(
					GTE(1),
					LTE(7),
				))).Validate(str)
				So(err, ShouldBeNil)
			})

			Convey("invalid", func() {
				str := normFields{}
				str.String.Scan("12345678")
				err := NormField("String", true,
					Length(All(
						GTE(1),
						LTE(7),
					))).Validate(str)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Decimal comparison", func() {
			d := field.Decimal{}
			d.Scan("4.50")
			err := LTE("5.00").Validate(d.Dec.String())
			So(err, ShouldBeNil)
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

		Convey("Errors", func() {
			Convey("Embedded", func() {
				fl := fieldList{
					Raws: []raw{
						{Int64: 100},
						{Int64: 120},
					},
				}

				err := Field("Raws", List(Field("Int64", LT(100)))).Validate(fl)
				ves, ok := err.(ValidationErrors)
				So(ok, ShouldBeTrue)
				So(len(ves), ShouldEqual, 2)
				So(ves[0].Location(), ShouldEqual, "Raws[0].Int64")
				So(ves[1].Location(), ShouldEqual, "Raws[1].Int64")
			})

			Convey("List in List", func() {
				ll := fieldListinList{
					List: []fieldList{
						{Raws: []raw{
							{String: "adam"},
						}},
						{Raws: []raw{
							{String: "adam", Int64: 110},
						}},
						{Raws: []raw{
							{String: "adam", Int64: 110},
						}},
					},
				}

				err := Field("List",
					All(
						List(
							Field("Raws", List(
								All(
									Field("String", Email),
									Field("Int64", LT(100)),
								),
							)),
						),
						Length(Equals(2)),
					),
				).Validate(ll)
				So(err, ShouldNotBeNil)
				ves, ok := err.(ValidationErrors)
				So(ok, ShouldBeTrue)
				So(len(ves), ShouldEqual, 6)
				So(ves[0].Location(), ShouldEqual, "List[0].Raws[0].String")
				So(ves[1].Location(), ShouldEqual, "List[1].Raws[0].String")
				So(ves[2].Location(), ShouldEqual, "List[1].Raws[0].Int64")
				So(ves[3].Location(), ShouldEqual, "List[2].Raws[0].String")
				So(ves[4].Location(), ShouldEqual, "List[2].Raws[0].Int64")
				So(ves[5].Location(), ShouldEqual, "List")
			})

			Convey("Error Messages", func() {
				Convey("Comparisons", func() {
					Convey("GT", func() {
						gt5 := GT(5)
						err := gt5.Validate(5)
						ve := err.(ValidationError)
						So(ve.Err, ShouldEqual, "property should be greater than 5")
					})

					Convey("GTE", func() {
						gte5 := GTE(5)
						err := gte5.Validate(4)
						ve := err.(ValidationError)
						So(ve.Err, ShouldEqual, "property should be greater or equal to 5")
					})

					Convey("LT", func() {
						lt5 := LT(5)
						err := lt5.Validate(5)
						ve := err.(ValidationError)
						So(ve.Err, ShouldEqual, "property should be less than 5")
					})

					Convey("LTE", func() {
						lte5 := LTE(5)
						err := lte5.Validate(6)
						ve := err.(ValidationError)
						So(ve.Err, ShouldEqual, "property should be less or equal to 5")
					})

					Convey("Length", func() {
						lengthlt5 := Length(LT(5))
						err := lengthlt5.Validate("hello")
						ve := err.(ValidationError)
						So(ve.Err, ShouldEqual, "length property should be less than 5")
					})
				})
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
