package norm

import (
	"fmt"
	"github.com/AlekSi/reflector"
	"github.com/picatic/go-api/norm/field"
)

func escapeFields(fields field.FieldNames) []string {
	var newFields []string = make([]string, len(fields))
	for i := 0; i < len(fields); i++ {
		newFields[i] = fmt.Sprintf("`%s`", fields[i].SnakeCase())
	}
	return newFields
}

func defaultFieldsEscaped(model Model, fields field.FieldNames) []string {
	if fields == nil {
		fields = ModelFields(model)
	}

	return escapeFields(fields)
}

// func fieldToDbMap(m Model) map[field.FieldName]string {
// 	fields := ModelFields(m)
// 	dbMap := make(map[field.FieldName]string)
// 	for i := 0; i < len(fields); i++ {
// 		dbMap[fields[i]] = dbr.NameMapping(string(fields[i]))
// 	}

// 	return dbMap
// }

func defaultUpdate(m Model, fields field.FieldNames) map[string]interface{} {
	kv := make(map[string]interface{})
	reflector.StructToMap(m, kv, "db")
	if fields == nil {
		return kv
	}
	fv := make(map[string]interface{})
	for _, k := range fields {
		if val, ok := kv[string(k)]; ok {
			fv[k.SnakeCase()] = val
		}
	}
	return fv
}
