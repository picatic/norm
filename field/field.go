package field

import (
	"database/sql"
	"database/sql/driver"
	"github.com/gocraft/dbr"
)

// Shadower Support for shadow fields. Allows us to determine if a field has been altered or not.
type Shadower interface {
	ShadowValue() (driver.Value, error)
	IsDirty() bool
}

// Name The name of a field on a model
type Name string

// SnakeCase Returns a field as SnakeCase
func (fn Name) SnakeCase() string {
	return dbr.NameMapping(string(fn))
}

// Names A set of Names
type Names []Name

// SnakeCase Return []string of snake_case field names for database map
func (fn Names) SnakeCase() []string {
	snakes := make([]string, len(fn))
	for i := 0; i < len(fn); i++ {
		snakes[i] = fn[i].SnakeCase()
	}
	return snakes
}

// Has determines if the name is in the slice
func (fn Names) Has(name Name) bool {
	for _, i := range fn {
		if i == name {
			return true
		}
	}
	return false
}

// Remove Returns a new Names with the names provided removed
func (fn Names) Remove(names Names) Names {
	if len(names) == 0 {
		return names
	}
	newNames := Names{}
	for _, i := range fn {
		if names.Has(i) == false {
			newNames = append(newNames, i)
		}
	}
	return newNames
}

// Field Implementation required to get the basic norm features for field mapping and dirty
type Field interface {
	sql.Scanner   // we require Scanner implementations
	driver.Valuer // our values stand and guard for thee
	Shadower      // we require Shadower
}

// compile time check
var _ []sql.Scanner = []sql.Scanner{
	&String{},
	&NullString{},
	&Time{},
	&NullTime{},
	&Int64{},
	&NullInt64{},
	&Bool{},
	&NullBool{},
}

var _ []Field = []Field{
	&String{},
	&NullString{},
	&Time{},
	&NullTime{},
	&Int64{},
	&NullInt64{},
	&Bool{},
	&NullBool{},
}
