package norm

import (
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

func TestModel(t *testing.T) {
	Convey("Model", t, func() {
		model := &MockModel{}
		model.Id.Scan("1")
		model.Name.Scan("Mock")

		Convey("ModelFields", func() {
			fields := ModelFields(model)
			So(fields, ShouldContain, field.FieldName("Id"))
			So(fields, ShouldContain, field.FieldName("Name"))
			So(len(fields), ShouldEqual, 2)
		})
	})
}
