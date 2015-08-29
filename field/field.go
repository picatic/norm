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
