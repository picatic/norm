package norm

import (
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/picatic/norm/field"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// Mock Model for testing
type MockModel struct {
	Id        field.NullString `json:"id",sql:"id"`
	FirstName field.NullString `json:"first_name",sql:"first_name"`
	Org       field.NullString `json:"org",sql:"org"`
	Ignore    string
}

type MockModelEmbedded struct {
	MockModel
	Created  field.Time
	Modified field.Time
}

type MockModelInterfaceEmbedded struct {
	Model
}

func (*MockModel) TableName() string {
	return "mocks"
}

func (*MockModel) IsNew() bool {
	return false
}

func (*MockModel) PrimaryKey() PrimaryKeyer {
	return NewSinglePrimaryKey(field.Name("Id"))
}

func (MockModel) Validators() []FieldValidator {
	validators := make([]FieldValidator, 1)
	validators[0] = NewFieldValidator(field.Name("FirstName"), "matches", MockFieldValidatorFunc, "Pete")
	return validators
}

type MockModelCustomPrimaryKey struct {
	MockModel
}

func (*MockModelCustomPrimaryKey) PrimaryKey() PrimaryKeyer {
	return NewCustomPrimaryKey(field.Names{"Id"}, func(pk PrimaryKeyer, model Model) (field.Names, error) {
		f, _ := ModelGetField(model, field.Name("Id"))
		f.Scan("abc-123-xyz-789")
		return field.Names{"Id"}, nil
	})
}

func TestModel(t *testing.T) {
	Convey("Model", t, func() {
		db, mock, _ := sqlmock.New()
		conn := NewConnection(db, "mock_db", nil)

		model := &MockModel{}
		model.Id.Scan("1")
		model.FirstName.Scan("Mock")

		modelWithEmbedded := &MockModelEmbedded{}

		Convey("ModelFields", func() {
			Convey("On Ptr to Struct", func() {
				fields := ModelFields(model)
				So(fields, ShouldContain, field.Name("Id"))
				So(fields, ShouldContain, field.Name("FirstName"))
				So(fields, ShouldContain, field.Name("Org"))
				So(len(fields), ShouldEqual, 3)
			})

			Convey("With embedded struct", func() {
				fields := ModelFields(modelWithEmbedded)
				So(fields, ShouldContain, field.Name("Id"))
				So(fields, ShouldContain, field.Name("FirstName"))
				So(fields, ShouldContain, field.Name("Org"))
				So(fields, ShouldContain, field.Name("Created"))
				So(fields, ShouldContain, field.Name("Modified"))
				So(len(fields), ShouldEqual, 5)
			})

			Convey("With embedded Model interface", func() {
				m := &MockModelInterfaceEmbedded{model}
				fields := ModelFields(m)
				So(fields, ShouldContain, field.Name("Id"))
				So(fields, ShouldContain, field.Name("FirstName"))
				So(fields, ShouldContain, field.Name("Org"))
				So(len(fields), ShouldEqual, 3)
			})

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

			Convey("Field from Embedded struct", func() {
				modelWithEmbedded.Id.Scan("12")
				idField, err := ModelGetField(modelWithEmbedded, "Id")
				So(err, ShouldBeNil)
				f, ok := idField.(*field.NullString)
				So(ok, ShouldBeTrue)
				So(f.String, ShouldEqual, "12")
			})

			Convey("Fields from embedded Model interface", func() {
				m := &MockModelInterfaceEmbedded{Model: model}
				rawModelField, err := ModelGetField(m, "Id")
				So(err, ShouldBeNil)

				f, ok := rawModelField.(*field.NullString)
				So(ok, ShouldBeTrue)
				So(f.String, ShouldEqual, "1")
			})
		})

		Convey("ModelTableName", func() {
			Convey("Builds table name with connection database prepended", func() {
				sess := conn.NewSession(nil)
				So(ModelTableName(sess, model), ShouldEqual, "mock_db.mocks")
			})
		})

		Convey("NewSelect", func() {

			Convey("Without fields", func() {
				mock.ExpectQuery("SELECT `id`, `first_name`, `org` FROM mock_db\\.mocks").WillReturnRows(sqlmock.NewRows([]string{"id", "first_name"}).FromCSVString("2,mocker"))
				err := NewSelect(conn.NewSession(nil), model, nil).LoadStruct(model)
				So(err, ShouldBeNil)
			})

			Convey("With fields", func() {
				mock.ExpectQuery("SELECT `id` FROM mock_db\\.mocks").WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("2"))
				err := NewSelect(conn.NewSession(nil), model, field.Names{"Id"}).LoadStruct(model)
				So(err, ShouldBeNil)
				So(model.Id.String, ShouldEqual, "2")
			})
		})

		Convey("NewInsert", func() {

			Convey("Without fields", func() {
				mock.ExpectExec("INSERT INTO `mock_db`\\.`mocks` \\(`first_name`,`org`\\) VALUES \\('Mock',NULL\\)").WillReturnResult(sqlmock.NewResult(2, 1))

				_, err := NewInsert(conn.NewSession(nil), model, nil).Record(model).Exec()
				So(err, ShouldBeNil)
			})

			Convey("With fields", func() {
				mock.ExpectExec("INSERT INTO `mock_db`\\.`mocks` \\(`first_name`\\) VALUES \\('Mock'\\)").WillReturnResult(sqlmock.NewResult(3, 1))
				_, err := NewInsert(conn.NewSession(nil), model, field.Names{"FirstName"}).Record(model).Exec()
				So(err, ShouldBeNil)
			})

			Convey("With Custom Primary Key", func() {
				modelCust := &MockModelCustomPrimaryKey{}
				modelCust.FirstName.Scan("Custom Key")
				mock.ExpectExec("INSERT INTO `mock_db`\\.`mocks` \\((`id`|,|`first_name`)+\\) VALUES \\(('abc-123-xyz-789'|,|'Custom Key')+\\)").WillReturnResult(sqlmock.NewResult(4, 1))
				_, err := NewInsert(conn.NewSession(nil), modelCust, field.Names{"FirstName"}).Record(modelCust).Exec()
				So(err, ShouldBeNil)
			})
		})

		Convey("NewUpdate", func() {

			Convey("Without fields", func() {
				mock.ExpectExec("UPDATE `mock_db`\\.`mocks` SET (`first_name` = 'Mock'|, |`org` = NULL)+ WHERE \\(id = '1'\\)").WillReturnResult(sqlmock.NewResult(0, 1))

				_, err := NewUpdate(conn.NewSession(nil), model, nil).Where("id = ?", model.Id.String).Exec()
				So(err, ShouldBeNil)
			})

			Convey("With fields", func() {
				mock.ExpectExec("UPDATE `mock_db`\\.`mocks` SET `first_name` = 'Mock' WHERE \\(id = '1'\\)").WillReturnResult(sqlmock.NewResult(0, 1))
				_, err := NewUpdate(conn.NewSession(nil), model, field.Names{"FirstName"}).Where("id = ?", model.Id.String).Exec()
				So(err, ShouldBeNil)
			})
		})

		Convey("ModelLoadMap", func() {
			dataMap := map[string]interface{}{
				"id":         "1234",
				"first_name": "James",
			}
			ModelLoadMap(model, dataMap)
			So(model.Id.String, ShouldEqual, "1234")
			So(model.FirstName.String, ShouldEqual, "James")
		})

		Convey("ModelChangedFields", func() {
			model := &MockModel{}
			model.Id.Scan("1")
			model.FirstName.Scan("James James James")
			Convey("No changed fields", func() {
				f, err := ModelDirtyFields(model)
				So(len(f), ShouldEqual, 0)
				So(err, ShouldBeNil)
			})

			Convey("Changed", func() {
				model.FirstName.Scan("Santa")
				f, err := ModelDirtyFields(model)
				So(len(f), ShouldEqual, 1)
				So(err, ShouldBeNil)
			})
		})

		Convey("ModelGetSetFields", func() {
			model := &MockModel{}

			Convey("No set fields", func() {
				f, err := ModelGetSetFields(model)
				So(len(f), ShouldEqual, 0)
				So(err, ShouldBeNil)
			})

			Convey("Changed", func() {
				model.Id.Scan("1")
				model.FirstName.Scan("James James James")
				f, err := ModelGetSetFields(model)
				So(len(f), ShouldEqual, 2)
				So(err, ShouldBeNil)
			})
		})

		Convey("ModelValidate", func() {

			Convey("With validation error", func() {
				err := ModelValidate(conn.NewSession(nil), model, nil)
				So(err, ShouldNotBeNil)
				So(err, ShouldHaveSameTypeAs, &ValidationErrors{})
				So(err.Error(), ShouldEqual, "Field: [FirstName] Alias: [matches] Message: Value [Mock] did not equal first argument [Pete]")
			})

			Convey("Without validation error", func() {
				model.FirstName.Scan("Pete")
				err := ModelValidate(conn.NewSession(nil), model, nil)
				So(err, ShouldBeNil)
			})

		})
	})
}
