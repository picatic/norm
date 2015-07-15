package norm

import (
	"errors"
	"github.com/gocraft/dbr"
	"github.com/picatic/go-api/norm/field"
	"reflect"
)

// All models have to implement this
type Model interface {
	TableName() string                    // table name within the database this model is associated to
	PrimaryKeyFieldName() field.FieldName // primary key for this model
	IsNew() bool                          // Is this model new or not
}

// Fetch list of fields on this model via reflection of fields that are from norm/field
func ModelFields(model Model) field.FieldNames {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() != reflect.Ptr {
		panic("Expected a Ptr")
	}

	if modelType.Elem().Kind() != reflect.Struct {
		panic("Expected struct")
	}
	modelValue := reflect.ValueOf(model).Elem()
	fieldType := reflect.TypeOf((*field.Field)(nil)).Elem()

	fields := make(field.FieldNames, 0)

	for i := 0; i < modelValue.NumField(); i++ {
		_field := modelValue.Field(i)
		if _field.CanAddr() == true && _field.Addr().Type().Implements(fieldType) == true {
			fields = append(fields, field.FieldName(modelType.Elem().Field(i).Name))
		}
	}

	return fields
}

// Get a field on a model by name
func ModelGetField(model Model, fieldName field.FieldName) (field.Field, error) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() != reflect.Ptr {
		panic("Expected a Ptr")
	}

	if modelType.Elem().Kind() != reflect.Struct {
		panic("Expected struct")
	}

	if _, ok := modelType.Elem().FieldByName(string(fieldName)); ok == true {
		modelValue := reflect.ValueOf(model).Elem().FieldByName(string(fieldName)).Addr().Interface()
		return modelValue.(field.Field), nil
	} else {
		return nil, errors.New("FieldName not found")
	}
}

//
// TODO: Would be nice to have the dbr.Session reliant code in a sub package...maybe.
// This is kind of an ActiveRecord/RemoteProxy/RecordGateway pattern
//

// NewSelect builds a select from the Model and Fields
// Selects all fields if no fields provided
func NewSelect(s *dbr.Session, m Model, fields field.FieldNames) *dbr.SelectBuilder {
	return s.Select(defaultFieldsEscaped(m, fields)...).From(m.TableName())
}

//NewUpdate builds an update from the Model and Fields
func NewUpdate(s *dbr.Session, m Model, f field.FieldNames) *dbr.UpdateBuilder {
	if f == nil {
		f = ModelFields(m)
	}
	setMap := defaultUpdate(m, f)
	return s.Update(m.TableName()).SetMap(setMap)
}

//NewInsert create an insert from the Model and Fields
func NewInsert(s *dbr.Session, m Model, fields field.FieldNames) *dbr.InsertBuilder {
	if fields == nil {
		fields = ModelFields(m)
	}
	return s.InsertInto(m.TableName()).Columns(fields.SnakeCase()...)
	// return s.InsertInto(m.TableName()).Columns(fields.SnakeCase()...)
}

//NewDelete creates a delete from the Model
func NewDelete(s *dbr.Session, m Model) *dbr.DeleteBuilder {
	return s.DeleteFrom(m.TableName())
}

// Save a model, calls appropriate Insert or Update based on Model.IsNew()
func ModelSave(dbrSess *dbr.Session, model Model, fields field.FieldNames) error {
	if model.IsNew() == true {
		return nil
	} else {
		return nil
	}
}
