package norm

import (
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gocraft/dbr"
	"github.com/picatic/go-api/norm/field"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// Mock Model for testing
type MockModel struct {
	Id     field.NullString
	Name   field.NullString
	Ignore string
}

func (MockModel) TableName() string {
	return "mocks"
}

func (MockModel) IsNew() bool {
	return false
}

func (MockModel) PrimaryKeyFieldName() field.FieldName {
	return "Id"
}

func TestModel(t *testing.T) {
	Convey("Model", t, func() {
		db, _ := sqlmock.New()
		dbrConn := dbr.NewConnection(db, nil)

		model := &MockModel{}
		model.Id.Scan("1")
		model.Name.Scan("Mock")

		Convey("ModelFields", func() {
			fields := ModelFields(model)
			So(fields, ShouldContain, field.FieldName("Id"))
			So(fields, ShouldContain, field.FieldName("Name"))
			So(len(fields), ShouldEqual, 2)
		})

		Convey("ModelGetField", func() {

			Convey("When field exists", func() {
				rawModelField, err := ModelGetField(model, "Id")
				So(err, ShouldBeNil)

				f, ok := rawModelField.(*field.NullString)
				So(ok, ShouldBeTrue)
				So(f.String, ShouldEqual, "1")
			})

			Convey("When field does not exist", func() {
				rawModelField, err := ModelGetField(model, "NotAField")
				So(rawModelField, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("NewSelect", func() {

			Convey("Without fields", func() {
				sqlmock.ExpectQuery("SELECT `Id`, `Name` FROM mocks").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).FromCSVString("2,mocker"))
				err := NewSelect(dbrConn.NewSession(nil), model, nil).LoadStruct(model)
				So(err, ShouldBeNil)
			})

			Convey("With fields", func() {
				sqlmock.ExpectQuery("SELECT `Id`, `Name` FROM mocks").WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("2"))
				err := NewSelect(dbrConn.NewSession(nil), model, nil).LoadStruct(model)
				So(err, ShouldBeNil)
			})
		})
	})
}
