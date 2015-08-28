package main

import (
	"github.com/picatic/norm/field"
)

type User struct {
	Id        field.NullInt64
	FirstName field.String
	LastName  field.String
	Email     field.String
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) PrimaryKeyFieldName() field.Name {
	return field.Name("Id")
}

func (u *User) IsNew() bool {
	return u.Id.Valid
}
