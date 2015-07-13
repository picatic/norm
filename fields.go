package norm

import (
	"errors"
	"github.com/gocraft/dbr"
	"database/sql"
)

type FieldName string

type FieldShadow interface{
	ShadowValue() (interface{}, error)
	IsDirty() bool
}

type Field interface {
	FieldShadow
}

// String Field
type StringField struct {
	String string
	shadow string
}

func (s *StringField) Scan(value interface{}) error {
	s.String = value.(string)

	return nil
}




// Nullable String
type NullString struct {
	dbr.NullString
	shadow dbr.NullString
	isShadowCreated OnceDone
}

func (ns *NullString) Scan(value interface{}) error {

	err := ns.NullString.Scan(value)
	if err != nil {
		return err
	}

	// load shadow on first scan only
	ns.isShadowCreated.Do(func(){
		_ = ns.shadow.Scan(value)
	})
	return nil
}

func (ns *NullString) Value() (interface{}, error) {
	if ns.Valid != true {
		return nil, nil
	}


	return ns.String, nil
}

func (ns *NullString) IsDirty() bool {
	if ns.Valid != ns.shadow.Valid || ns.String != ns.shadow.String {
		return true
	}
	return false
}

func (ns *NullString) ShadowValue() (interface{}, error) {
	if ns.isShadowCreated.Done() {
		return ns.shadow.Value()
	}
	return nil, errors.New("Shadow Wasn't Created")
}

// compile time check
var _ sql.Scanner = &NullString{}