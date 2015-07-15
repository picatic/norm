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
	Id        field.NullString `json:"id",sql:"id"`
	FirstName field.NullString `json:"first_name",sql:"first_name"`
	Ignore    string
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
		model.FirstName.Scan("Mock")

		Convey("ModelFields", func() {
			fields := ModelFields(model)
			So(fields, ShouldContain, field.FieldName("Id"))
			So(fields, ShouldContain, field.FieldName("FirstName"))
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
				sqlmock.ExpectQuery("SELECT `id`, `first_name` FROM mocks").WillReturnRows(sqlmock.NewRows([]string{"id", "first_name"}).FromCSVString("2,mocker"))
				err := NewSelect(dbrConn.NewSession(nil), model, nil).LoadStruct(model)
				So(err, ShouldBeNil)
			})

			Convey("With fields", func() {
				sqlmock.ExpectQuery("SELECT `id` FROM mocks").WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("2"))
				err := NewSelect(dbrConn.NewSession(nil), model, field.FieldNames{"Id"}).LoadStruct(model)
				So(err, ShouldBeNil)
			})
		})

		Convey("NewInsert", func() {

			Convey("Without fields", func() {
				sqlmock.ExpectExec("INSERT INTO mocks \\(`id`,`first_name`\\) VALUES \\('1','Mock'\\)").WillReturnResult(sqlmock.NewResult(2, 1))

				_, err := NewInsert(dbrConn.NewSession(nil), model, nil).Record(model).Exec()
				// log.Printf("%+v\n", ib.Cols)
				So(err, ShouldBeNil)
			})

			Convey("With fields", func() {
				sqlmock.ExpectExec("INSERT INTO mocks \\(`first_name`\\) VALUES \\('Mock'\\)").WillReturnResult(sqlmock.NewResult(3, 1))
				_, err := NewInsert(dbrConn.NewSession(nil), model, field.FieldNames{"FirstName"}).Record(model).Exec()
				So(err, ShouldBeNil)
			})
		})
	})
}
