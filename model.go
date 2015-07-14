package norm

import (
	"errors"
	"github.com/gocraft/dbr"
	"github.com/picatic/go-api/norm/field"
)

// All models have to implement this
type Model interface {
	TableName() string                         // table name within the database this model is associated to
	GetField(name field.FieldName) interface{} // return a norm.field by name
	IsNew() bool                               // Is this model new or not
}

// Fetch list of fields on this model via reflection of fields that are from norm/field
func ModelFields(model Model) []field.FieldName {
	// reflect ?
	panic(errors.New("NotImplemented"))
	return nil
}

//
// TODO: Would be nice to have the dbr.Session reliant code in a sub package...maybe.
// This is kind of an ActiveRecord/RemoteProxy/RecordGateway pattern
//
func ModelLoad(dbrSess *dbr.Session, model Model) error {
	return errors.New("NotImplemented")
}

func ModelLoadMany(dbrSess *dbr.Session, model []Model) error {
	return errors.New("NotImplemented")
}

func ModelSave(dbrSess *dbr.Session, model Model) error {
	if model.IsNew() == true {
		return modelCreate(dbrSess, model, model)
	} else {
		return errors.New("NotImplemented")
	}
}

func modelCreate(dbrSess *dbr.Session, model Model, fields []field.FieldName) error {
	//	return NewInsert(dbrSess, model, model.Fields()).Record(model).Exec()
	return errors.New("NotImplemented")
}

func modelUpdate() {

}

func ModelSaveFields(dbrSess *dbr.Session, model Model, fields []field.FieldName) error {
	return errors.New("NotImplemented")
}
