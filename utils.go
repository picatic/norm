package norm

import (
	"fmt"
	_ "github.com/AlekSi/reflector"
	"github.com/picatic/go-api/norm/field"
)

func escapeFields(fields field.FieldNames) []string {
	var newFields []string = make([]string, len(fields))
	for i := 0; i < len(fields); i++ {
		newFields[i] = fmt.Sprintf("`%s`", string(fields[i]))
	}
	return newFields
}

func defaultFieldsEscaped(model Model, fields field.FieldNames) []string {
	if fields == nil {
		fields = ModelFields(model)
	}

	return escapeFields(fields)
}

// func defaultUpdate(m Model, fields field.FieldNames) map[string]interface{} {
// 	kv := make(map[string]interface{})
// 	reflector.StructToMap(m, kv, "db")
// 	if fields == nil {
// 		return kv
// 	}
// 	fv := make(map[string]interface{})
// 	for _, k := range fields.Fields() {
// 		if val, ok := kv[k]; ok {
// 			fv[k] = val
// 		}
// 	}
// 	return fv
// }
