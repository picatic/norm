package main

import (
	"github.com/picatic/norm/field"
)

type User struct {
	Id        field.NullInt64  `json:"id"`
	FirstName field.NullString `json:"first_name"`
	LastName  field.String     `json:"last_name"`
	Email     field.String     `json:"email"`
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
