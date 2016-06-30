package field

import (
	"database/sql"
	"database/sql/driver"
	"unicode"
)

// Shadower Support for shadow fields. Allows us to determine if a field has been altered or not.
type Shadower interface {
	ShadowValue() (driver.Value, error)
	IsDirty() bool
}

// Name The name of a field on a model
type Name string

// NewNameFromSnakeCase snake_case to CamelCase
func NewNameFromSnakeCase(name string) Name {
	var newstr []rune
	nextCap := true

	for _, chr := range name {
		if nextCap == true && chr != '_' {
			chr += ('A' - 'a')
			newstr = append(newstr, chr)
			nextCap = false
		} else if chr == '_' {
			nextCap = true
		} else {
			newstr = append(newstr, chr)
		}
	}
	return Name(string(newstr))
}

// SnakeCase Returns a field as SnakeCase
func (fn Name) SnakeCase() string {
	var newstr []rune
	lastLower := false
	for _, chr := range fn {
		if unicode.IsUpper(chr) {
			if lastLower {
				newstr = append(newstr, '_')
			}
			newstr = append(newstr, unicode.ToLower(chr))
			lastLower = false
		} else if unicode.IsLower(chr) {
			newstr = append(newstr, chr)
			lastLower = true
		} else {
			newstr = append(newstr, chr)
			lastLower = false
		}
	}
	return string(newstr)
}

// implements Locationer from validate
func (fn Name) Location() string {
	return string(fn)
}

// Names A set of Names
type Names []Name

// NewNamesFromString create Names from []string
func NewNamesFromString(strs []string) Names {
	names := make(Names, len(strs), len(strs))
	for i, v := range strs {
		names[i] = Name(string(v))
	}
	return names
}

// NewNamesFromSnakeCase create Names with snake_case []string
func NewNamesFromSnakeCase(snakes []string) Names {
	names := make(Names, len(snakes), len(snakes))
	for i, v := range snakes {
		names[i] = NewNameFromSnakeCase(v)
	}
	return names
}

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
		return fn
	}
	newNames := Names{}
	for _, i := range fn {
		if names.Has(i) == false {
			newNames = append(newNames, i)
		}
	}
	return newNames
}

// Intersect returns the intersection of the two Names
func (fn Names) Intersect(names Names) Names {
	var union Names
	for _, field := range names {
		if fn.Has(field) {
			union = append(union, field)
		}
	}
	return union
}

// Add returns a new Names with the names provided added, no duplicates
func (fn Names) Add(names Names) Names {
	if len(names) == 0 {
		return fn
	}
	newNames := fn
	for _, i := range names {
		if fn.Has(i) == false {
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
	IsSet() bool  // True if Scan has been called
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
