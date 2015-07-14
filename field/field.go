package field

import (
"errors"
"github.com/picatic/go-api/norm"
"github.com/gocraft/dbr"
"database/sql"
)

type FieldName string

type FieldShadow interface {
ShadowValue() (interface{}, error)
IsDirty() bool
}

type Field interface {
FieldShadow
}

//
// String
//
type String struct {
String string
shadow string
}

func (s *String) Scan(value interface{}) error {
s.String = value.(string)

return nil
}

var _ sql.Scanner = &String{}
var _ Field = &String{}


//
// NullString
//
type NullString struct {
dbr.NullString
shadow          dbr.NullString
isShadowCreated norm.OnceDone
}

func (ns *NullString) Scan(value interface{}) error {

err := ns.NullString.Scan(value)
if err != nil {
return err
}

// load shadow on first scan only
ns.isShadowCreated.Do(func() {
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
return  ns.Valid != ns.shadow.Valid || ns.String != ns.shadow.String
}

func (ns *NullString) ShadowValue() (interface{}, error) {
if ns.isShadowCreated.Done() {
return ns.shadow.Value()
}

return nil, errors.New("Shadow Wasn't Created")
}

// compile time check
var _ sql.Scanner = &NullString{}
var _ Field = &NullString{}