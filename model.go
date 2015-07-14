package norm

import (
	"github.com/gocraft/dbr"
	"errors"
)

type FieldName string

type Fields interface {
	Fields() []FieldName // go field camel case name
	GetField(name FieldName) interface{} // return a field
}

type Table interface {
	TableName() string
}

type Model interface {
	Table
	Fields
	IsNew() bool
}

func ModelFields(model Model) []FieldName {
	// reflect ?
	panic(errors.New("NotImplemented"))
	return nil
}

//
// TODO: Would be nice to have the dbr.Session reliant code in a sub package...maybe.
// This is kind of an ActiveRecord/RemoteProxy/RecordGateway pattern
//
func ModelLoad(dbrSess *dbr.Session, model Model) (error ) {
	return errors.New("NotImplemented")
}

func ModelLoadMany(dbrSess *dbr.Session, model []Model) (error) {
	return errors.New("NotImplemented")
}

func ModelSave(dbrSess *dbr.Session, model Model) (error) {
	if model.IsNew() == true {
		return modelCreate(dbrSess, model, model)
	} else {
		return errors.New("NotImplemented")
	}
}

func modelCreate(dbrSess *dbr.Session, model Model, fields Fields) (error) {
//	return NewInsert(dbrSess, model, model.Fields()).Record(model).Exec()
	return errors.New("NotImplemented")
}

func modelUpdate() {

}

func ModelSaveFields(dbrSess *dbr.Session, model Model, fields []Fields) (error) {
	return errors.New("NotImplemented")
}