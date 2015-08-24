package norm

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gocraft/dbr"
	"github.com/picatic/norm/field"
	"reflect"
)

// Models use this interface to communicate some basic information to help with building queries
// and determining the state of the model.
type Model interface {
	TableName() string               // table name within the database this model is associated to
	PrimaryKeyFieldName() field.Name // primary key for this model
	IsNew() bool                     // Is this model new or not
}

// Fetch list of fields on this model via reflection of fields that are from norm/field
// If model fails to be a Ptr to a Struct we return an error
func ModelFields(model Model) field.Names {
	fieldType := reflect.TypeOf((*field.Field)(nil)).Elem()

	modelType := reflect.TypeOf(model)

	if modelType.Kind() != reflect.Ptr {
		panic("Expected Model to be a Ptr")
	}

	if modelType.Elem().Kind() != reflect.Struct {
		panic("Expected Model to be a Ptr to a Struct")
	}

	modelValue := reflect.ValueOf(model).Elem()

	fields := make(field.Names, 0)

	for i := 0; i < modelValue.NumField(); i++ {
		_field := modelValue.Field(i)
		if _field.CanAddr() == true && _field.Addr().Type().Implements(fieldType) == true {
			fields = append(fields, field.Name(modelType.Elem().Field(i).Name))
		}
	}

	return fields
}

// Get a field on a model by name
func ModelGetField(model Model, fieldName field.Name) (field.Field, error) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() != reflect.Ptr {
		panic("Expected Model to be a Ptr")
	}

	if modelType.Elem().Kind() != reflect.Struct {
		panic("Expected Model to be a Ptr to Struct")
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
func NewSelect(s *dbr.Session, m Model, fields field.Names) *dbr.SelectBuilder {
	return s.Select(defaultFieldsEscaped(m, fields)...).From(m.TableName())
}

//NewUpdate builds an update from the Model and Fields
func NewUpdate(s *dbr.Session, m Model, f field.Names) *dbr.UpdateBuilder {
	if f == nil {
		f = ModelFields(m)
	}
	setMap := defaultUpdate(m, f)
	return s.Update(m.TableName()).SetMap(setMap)
}

//NewInsert create an insert from the Model and Fields
func NewInsert(s *dbr.Session, m Model, fields field.Names) *dbr.InsertBuilder {
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
func ModelSave(dbrSess *dbr.Session, model Model, fields field.Names) (sql.Result, error) {
	if model.IsNew() == true {
		return nil, errors.New("ModelSave for when IsNew Not Implement")
	} else {
		primaryFieldName := model.PrimaryKeyFieldName()
		idField, err := ModelGetField(model, primaryFieldName)
		if err != nil {
			return nil, err
		}
		id, err := idField.Value()
		if err != nil {
			return nil, err
		}
		return NewUpdate(dbrSess, model, fields).Where(fmt.Sprintf("`%s`=?", primaryFieldName.SnakeCase()), id).Exec()
	}
}

// Validate fields provided on model, if no fields validate all fields
func ModelValidate(model Model, fields field.Names) chan error {
	err := make(chan error, 1)

	go func() {
		if fields == nil {
			fields = ModelFields(model)
		}
		errors := Validators.Validate(model, fields)
		if errors != nil {
			for _, validationError := range errors.Errors {
				err <- validationError
			}
		}
		err <- nil
	}()
	return err
}
