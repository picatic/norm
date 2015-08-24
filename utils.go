package norm

import (
	"fmt"
	"github.com/AlekSi/reflector"
	"github.com/picatic/norm/field"
)

// escape fields for queries
func escapeFields(fields field.Names) []string {
	var newFields []string = make([]string, len(fields))
	for i := 0; i < len(fields); i++ {
		newFields[i] = fmt.Sprintf("`%s`", fields[i].SnakeCase())
	}
	return newFields
}

// Get the fieldNames as a []string escaped field names
func defaultFieldsEscaped(model Model, fields field.Names) []string {
	if fields == nil {
		fields = ModelFields(model)
	}

	return escapeFields(fields)
}

// func fieldToDbMap(m Model) map[field.Name]string {
// 	fields := ModelFields(m)
// 	dbMap := make(map[field.Name]string)
// 	for i := 0; i < len(fields); i++ {
// 		dbMap[fields[i]] = dbr.NameMapping(string(fields[i]))
// 	}

// 	return dbMap
// }

// Create a map of strings and values from the model to work with dbr's interfaces
func defaultUpdate(m Model, fields field.Names) map[string]interface{} {
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
