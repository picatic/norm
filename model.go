package norm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/gocraft/dbr"
	"github.com/picatic/norm/field"
	"github.com/picatic/norm/validate"
	"github.com/picatic/norm/atomiccache"
)

var (
	modelFieldsCache atomiccache.Cache
	fieldType       = reflect.TypeOf((*field.Field)(nil)).Elem()
	modelType       = reflect.TypeOf((*Model)(nil)).Elem()
	NameNotFoundErr = errors.New("Name not found")
)

// Model This interface provides basic information to help with building queries and behaviours in dbr.
type Model interface {
	// TableName is the name used in the name used in database.
	//
	// Typically implemented to return in storage table name
	//
	//	// TableName returns database table name
	//	func (u *User) TableName() string {
	//		return "norm_users"
	//	}
	TableName() string

	// PrimaryKey returns a PrimaryKeyer which can be one of a few provided implementations
	//
	//	// PrimaryKey returns SinglePrimaryKey
	//	func (u *User) PrimaryKey() PrimaryKeyer {
	//		return norm.NewSinglePrimaryKey(field.Name("Id"))
	//	}
	//
	//	// PrimaryKey returns Composite Key as MultiplePrimaryKey
	//	func (u *User) PrimaryKey() PrimaryKeyer {
	//		return norm.NewMultiplePrimaryKey(field.Names{"OrgId", "AccountId"})
	//	}
	//
	//	// PrimaryKey returns CustomPrimaryKey example for generating a uuid
	//	func (t OauthAccessToken) PrimaryKey() norm.PrimaryKeyer {
	//		return norm.NewCustomPrimaryKey(Names{"Id"}, func(pk norm.PrimaryKeyer, model norm.Model) (Names, error) {
	//			f, _ := norm.ModelGetField(model, "Id")
	//			f.Scan(newUuid())
	//			return Names{"Id"}, nil
	//		})
	//	}
	//
	//  Must not be called more than once and should only be done by norm.
	PrimaryKey() PrimaryKeyer

	// IsNew returns if the the model does not already exist in storage.
	// The Primarykey on the model should not be Valid.
	//
	// Typically this can be implemented as
	//
	//	// IsNew indicates that model instance does not already exist.
	//	func (u *User) IsNew() bool {
	//		return !u.Id.Valid
	//	}
	IsNew() bool
}

// ModelFields Fetch list of fields on this model via reflection of fields that are from norm/field
// If model fails to be a Ptr to a Struct we return an error
func ModelFields(model Model) field.Names {
	modelType := reflect.TypeOf(model)

	return modelFieldsCache.Get(modelType, func() interface{} {
		if modelType.Kind() != reflect.Ptr {
			panic("Expected Model to be a Ptr")
		}

		if modelType.Elem().Kind() != reflect.Struct {
			panic("Expected Model to be a Ptr to a Struct")
		}

		return modelFields(model)
	}).(field.Names)
}

func modelFields(model interface{}) field.Names {
	fields := make(field.Names, 0)

	// Value of model
	ifv := reflect.ValueOf(model)
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
	}

	// Type of Model
	itf := reflect.TypeOf(model)
	if itf.Kind() == reflect.Ptr {
		itf = itf.Elem()
	}

	for i := 0; i < itf.NumField(); i++ {
		v := ifv.Field(i)

		// Straight up struct field of type field.Field
		if v.CanAddr() == true && v.Addr().Type().Implements(fieldType) == true {
			fields = append(fields, field.Name(itf.Field(i).Name))
		} else {
			t := itf.Field(i)
			// Embedded Struct with potential field.Field fields
			if t.Anonymous == true && v.CanAddr() == true && v.Kind() == reflect.Struct {
				fields = append(fields, modelFields(v.Addr().Interface())...)
			} else if t.Anonymous == true && v.CanAddr() == true && v.Kind() == reflect.Interface { // Embedded Model interface
				fields = append(fields, modelFields(v.Elem().Interface())...)
			}
		}
	}

	return fields
}

// ModelGetField Get a field on a model by name.
// This function  will just returns field.Field or NameNotFoundErr error
func ModelGetField(model Model, fieldName field.Name) (field.Field, error) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() != reflect.Ptr {
		panic("Expected Model to be a Ptr")
	}

	if modelType.Elem().Kind() != reflect.Struct {
		panic("Expected Model to be a Ptr to Struct")
	}
	return modelGetField(model, fieldName)
}

func modelGetField(model interface{}, fieldName field.Name) (field.Field, error) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	if _, ok := modelType.FieldByName(string(fieldName)); ok == true {
		modelValue := reflect.ValueOf(model).Elem().FieldByName(string(fieldName)).Addr().Interface()
		return modelValue.(field.Field), nil
	} else {
		modelValue := reflect.ValueOf(model)

		if modelValue.Kind() == reflect.Ptr {
			modelValue = modelValue.Elem()
		}
		for i := 0; i < modelType.NumField(); i++ {
			t := modelType.Field(i)
			v := modelValue.Field(i)

			if t.Anonymous == true && v.CanAddr() == true && v.Kind() == reflect.Interface {
				if m, ok := v.Elem().Interface().(Model); ok == true {
					if f, err := modelGetField(m, fieldName); err == nil {
						return f, nil
					}
				}

			}
		}
	}
	return nil, NameNotFoundErr
}

// ModelGetSetFields is named poorly but returns all the fields on a model that have been set.
// For a field to be set, it must of been successfully called with Scan at least once.
func ModelGetSetFields(model Model) (field.Names, error) {
	fields := ModelFields(model)

	var setFields field.Names

	for _, field := range fields {
		modelField, err := ModelGetField(model, field)
		if err != nil {
			return nil, err
		}
		if modelField.IsSet() {
			setFields = append(setFields, field)
		}
	}
	return setFields, nil
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
	// TODO: handle composite primary keys
	if len(model.PrimaryKey().Fields()) > 1 {
		panic("ModelSave Composite Primary Keys not yet implemented")
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

// ModelSetDefaults load a map into a model to set default values
func ModelSetDefaults(model Model, defaults map[string]interface{}) error {
	for k, v := range defaults {
		modelField, err := ModelGetField(model, field.NewNameFromSnakeCase(k))
		if err != nil {
			continue
		}
		if !modelField.IsSet() {
			err = modelField.Scan(v)
			if err != nil {
				return err
			}
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
	if validator, ok := model.(validate.ModelValidator); ok {
		return validator.Validator().Fields(fields).Validate(model)
	}

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
	err := validators.Validate(sess, model, fields)
	if err == nil { // a ptr to a struct cast into an interface is no longer really nil
		return nil
	}
	return err
}

//MapFields makes a copy from a origin model to a destination model fields with the same name or
//provided by a mapping table
func MapFields(originModel Model, destinationModel Model, mappings map[field.Name]field.Name) {
	originFieldNames := ModelFields(originModel)
	for _, originFieldName := range originFieldNames {
		var (
			destinationFieldName field.Name
			ok                   bool
		)

		if destinationFieldName, ok = mappings[originFieldName]; !ok {
			destinationFieldName = originFieldName
		}

		originField, err := ModelGetField(originModel, originFieldName)
		if err == NameNotFoundErr {
			continue
		}

		destinationField, err := ModelGetField(destinationModel, destinationFieldName)
		if err == NameNotFoundErr {
			continue
		}
		destinationField.Scan(originField)
	}
}

//MapSetFields makes a copy from a origin model to a destination model fields that are set with the same name or
//provided by a mapping table
func MapSetFields(originModel Model, destinationModel Model, mappings map[field.Name]field.Name) {
	originFieldNames := ModelFields(originModel)
	for _, originFieldName := range originFieldNames {
		var (
			destinationFieldName field.Name
			ok                   bool
		)

		if destinationFieldName, ok = mappings[originFieldName]; !ok {
			destinationFieldName = originFieldName
		}

		originField, err := ModelGetField(originModel, originFieldName)
		if err == NameNotFoundErr {
			continue
		}

		if originField.IsSet() {
			destinationField, err := ModelGetField(destinationModel, destinationFieldName)
			if err == NameNotFoundErr {
				continue
			}
			destinationField.Scan(originField)
		}
	}
}

func getAllFields(model interface{}) map[field.Name]reflect.Value {
	fields := map[field.Name]reflect.Value{}

	// Value of model
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	// Type of Model
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	for i := 0; i < modelType.NumField(); i++ {
		v := modelValue.Field(i)

		// Straight up struct field of type field.Field
		if v.CanAddr() == true && v.Addr().Type().Implements(fieldType) == true {
			fieldName := field.Name(modelType.Field(i).Name)
			fields[fieldName] = modelValue.Field(i)
		} else if v.Kind() == reflect.Struct {
			//t := modelType.Field(i)
			// Embedded Struct with potential field.Field fields
			fieldName := field.Name(modelType.Field(i).Name)
			fields[fieldName] = modelValue.Field(i)
		}
	}

	return fields
}

//ModelToMap return a map with Model fields
func ModelToMap(model interface{}, fields field.Names, mapEmbedded bool) (map[string]interface{}, error) {
	return modelToMap(model, fields, mapEmbedded)
}

func modelToMap(model interface{}, fields field.Names, mapEmbedded bool) (map[string]interface{}, error) {
	resultMap := map[string]interface{}{}
	allFields := getAllFields(model)

	if fields == nil {
		allFieldsKeys := reflect.ValueOf(allFields).MapKeys()
		fields = make(field.Names, len(allFieldsKeys))
		for i, key := range allFieldsKeys {
			fields[i] = key.Interface().(field.Name)
		}
	}

	for _, fieldName := range fields {
		fieldValue, ok := allFields[fieldName]
		if !ok {
			return nil, fmt.Errorf("Field name '%s' not found", fieldName)
		}
		//if it's a norm fields return the value
		if fieldValue.CanSet() == true && fieldValue.CanAddr() == true &&
			fieldValue.Addr().Type().Implements(fieldType) == true {
			field := fieldValue.Addr().Interface().(field.Field)
			value, err := field.Value()
			if err != nil {
				return nil, err
			}
			resultMap[fieldName.SnakeCase()] = value
			//if it's a embedded struct/model return the object map recursively
		} else if fieldValue.CanSet() == true && fieldValue.Kind() == reflect.Struct && mapEmbedded {
			var valueInterface interface{}
			if fieldValue.Kind() == reflect.Ptr {
				valueInterface = fieldValue.Interface()
			} else {
				valueInterface = fieldValue.Addr().Interface()
			}
			innerMap, err := modelToMap(valueInterface, nil, mapEmbedded)
			if err != nil {
				return nil, err
			}
			resultMap[fieldName.SnakeCase()] = innerMap
		}

	}

	return resultMap, nil
}
