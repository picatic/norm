package norm

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gocraft/dbr"
	"github.com/picatic/norm/field"
	"reflect"
)

// Model This interface provides basic information to help with building queries and behaviours in dbr.
type Model interface {
	TableName() string        // table name within the database this model is associated to
	PrimaryKey() PrimaryKeyer // primary key for this model
	IsNew() bool              // Is this model new or not
}

// ModelFields Fetch list of fields on this model via reflection of fields that are from norm/field
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

// ModelGetField Get a field on a model by name
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
	}
	return nil, errors.New("Name not found")
}

// ModelTableName get the complete table name including the database
func ModelTableName(s Session, m Model) string {
	return fmt.Sprintf("%s.%s", s.Connection().Database(), m.TableName())
}

//
// TODO: Would be nice to have the Session reliant code in a sub package...maybe.
// This is kind of an ActiveRecord/RemoteProxy/RecordGateway pattern
//

// NewSelect builds a select from the Model and Fields
// Selects all fields if no fields provided
func NewSelect(s Session, m Model, fields field.Names) *dbr.SelectBuilder {
	return s.Select(defaultFieldsEscaped(m, fields)...).From(ModelTableName(s, m))
}

// NewUpdate builds an update from the Model and Fields
func NewUpdate(s Session, m Model, fields field.Names) *dbr.UpdateBuilder {
	if fields == nil {
		fields = ModelFields(m)
	}
	fields = fields.Remove(m.PrimaryKey().Fields())
	setMap := defaultUpdate(m, fields)
	return s.Update(ModelTableName(s, m)).SetMap(setMap)
}

// NewInsert create an insert from the Model and Fields
func NewInsert(s Session, m Model, fields field.Names) *dbr.InsertBuilder {
	if fields == nil {
		fields = ModelFields(m)
	}
	pk := m.PrimaryKey()
	fields = fields.Remove(pk.Fields())
	// TODO do not eat this error
	setFields, _ := pk.Generator(m)
	fields = fields.Add(setFields)
	return s.InsertInto(ModelTableName(s, m)).Columns(fields.SnakeCase()...)
}

// NewDelete creates a delete from the Model
func NewDelete(s Session, m Model) *dbr.DeleteBuilder {
	return s.DeleteFrom(ModelTableName(s, m))
}

// ModelSave Save a model, calls appropriate Insert or Update based on Model.IsNew()
func ModelSave(dbrSess Session, model Model, fields field.Names) (sql.Result, error) {
	if model.IsNew() == true {
		return nil, errors.New("ModelSave for when IsNew Not Implement")
	}
	primaryFieldName := model.PrimaryKey().Fields()[0]

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

// ModelLoadMap load a map into a model
func ModelLoadMap(model Model, data map[string]interface{}) error {
	for k, v := range data {
		modelField, err := ModelGetField(model, field.NewNameFromSnakeCase(k))
		if err != nil {
			continue
		}

		err = modelField.Scan(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// ModelDirtyFields return Names of the fields that are dirty
func ModelDirtyFields(model Model) (field.Names, error) {
	dirtyFields := make(field.Names, 0)
	fields := ModelFields(model)

	for _, fieldName := range fields {
		mf, err := ModelGetField(model, fieldName)
		if err != nil {
			return nil, err
		}
		if mf.IsDirty() == true {
			dirtyFields = append(dirtyFields, fieldName)
		}
	}
	return dirtyFields, nil
}

// ModelValidate fields provided on model, if no fields validate all fields
func ModelValidate(sess Session, model Model, fields field.Names) error {
	validators := sess.Connection().ValidatorCache()
	if fields == nil {
		fields = ModelFields(model)
	}
	if len(validators.Get(model)) == 0 {
		if vm, ok := model.(ModelValidators); ok == true {
			validators.Set(model, vm.Validators())
		} else {
			return nil
		}
	}
	return validators.Validate(sess, model, fields)
}
