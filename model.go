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
	typeOfField := reflect.TypeOf((*field.Field)(nil)).Elem()
	r := reflect.TypeOf(model).Elem()
	fields := make([]field.FieldName, 1)
	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)
		if f.Type.Implements(typeOfField) == true {
			fields = append(fields, field.FieldName(f.Name))
		}
	}
	return fields
}

// Get a field on a model by name
func ModelGetField(model Model, field field.FieldName) (interface{}, error) {
	panic(errors.New("NotImplemented"))
	return nil, nil
}

//
// TODO: Would be nice to have the dbr.Session reliant code in a sub package...maybe.
// This is kind of an ActiveRecord/RemoteProxy/RecordGateway pattern
//

// NewSelect builds a select from the Model and Fields
// Selects all fields if no fields provided
func NewSelect(s *dbr.Session, m Model, f field.FieldNames) *dbr.SelectBuilder {
	return s.Select(defaultFieldsEscaped(m, f)...).From(m.TableName())
}

// load a model from a SelectBuilder
func ModelLoad(dbrSess *dbr.Session, model Model) error {
	return errors.New("NotImplemented")
}

func ModelLoadMany(dbrSess *dbr.Session, model []Model) error {
	return errors.New("NotImplemented")
}

//NewUpdate builds an update from the Model and Fields
func NewUpdate(s *dbr.Session, m Model, f field.FieldNames) *dbr.UpdateBuilder {
	panic("NotImplemented")
	//return s.Update(m.TableName()).SetMap(defaultUpdate(m, f))
	return nil
}

//NewInsert create an insert from the Model and Fields
func NewInsert(s *dbr.Session, m Model, f field.FieldNames) *dbr.InsertBuilder {
	return s.InsertInto(m.TableName()).Columns(defaultFieldsEscaped(m, f)...)
}

//NewDelete creates a delete from the Model
func NewDelete(s *dbr.Session, m Model) *dbr.DeleteBuilder {
	return s.DeleteFrom(m.TableName())
}

// Save a model
func ModelSave(dbrSess *dbr.Session, model Model) error {
	if model.IsNew() == true {
		return modelCreate(dbrSess, model, ModelFields(model))
	} else {
		return modelUpdate(dbrSess, model, ModelFields(model))
	}
}

func modelCreate(dbrSess *dbr.Session, model Model, fields []field.FieldName) error {
	//	return NewInsert(dbrSess, model, model.Fields()).Record(model).Exec()
	return errors.New("NotImplemented")
}

func modelUpdate(dbrSess *dbr.Session, model Model, fields []field.FieldName) error {
	return errors.New("NotImplemented")
}

// Save specific fields on a model
func ModelSaveFields(dbrSess *dbr.Session, model Model, fields []field.FieldName) error {
	if model.IsNew() == true {
		return modelCreate(dbrSess, model, fields)
	} else {
		return modelUpdate(dbrSess, model, fields)
	}
}
