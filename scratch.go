package norm

import (
	"github.com/gocraft/dbr"
	"github.com/picatic/go-api/norm/field"
	valid "github.com/asaskevich/govalidator"
)


type TimeField struct {
	dbr.NullTime
	Shadow dbr.NullTime
	Rules []Validator
}

type UserInterface interface {
	Name() field.String
	SetName(interface{})
}

type User struct {
	Name  field.String `json:"name",db:name`
	Title field.String `json:"title",db:title`
	Created TimeField `json:"created",db:created_at`
}

func (u *User) Fields() []FieldName {
	return []FieldName{"Name","Title","Created"}
}

//user := &User{}
//user.Name.Scan("James")

//
type StructValidator map[string][]Validator

func (sv *StructValidator) ValidateFields(s interface{}, fields []string) (bool, error) {
	for x := range fields {
		println(x)
		// validate desired fields
	}
	return true, nil
}

func ValidateStruct(s interface{}, fields []string) (bool, error) {
	//validate struct
	return true, nil
}

//var UserValidator StructValidator = map[string][]Validator {"foo": {Validator{}}}

type Validator interface {
	Validate(interface{}) (bool, error)
}

type ValidatorFunc func(interface{})(bool, error)

func (vf ValidatorFunc)Validate(field interface{}) (bool, error) {
	return vf(field)
}


var IsAlpha ValidatorFunc = func(f interface{}) (bool, error) {
	return valid.IsAlpha(f.(string)), nil }

var IsLength3To25 ValidatorFunc = func(f interface{}) (bool, error) {
	return valid.IsByteLength(f.(string), 3, 25), nil
}

func NewUser() User {
	u := User{}


//	u.Name.Rules = append(u.Name.Rules,
//		IsAlpha,
//		IsLength3To25)
	return u
}





