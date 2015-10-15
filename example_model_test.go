package norm_test

import (
	"encoding/json"
	"fmt"
	"github.com/picatic/norm"
	"github.com/picatic/norm/field"
	"os"
)

type User struct {
	Id        field.Int64  `json:"id",sql:"id"`
	FirstName field.String `json:"first_name"`
	LastName  field.String `json:"last_name"`
	Email     field.String `json:"email"`
}

func (u *User) TableName() string {
	return "norm_users"
}

func (u *User) PrimaryKey() norm.PrimaryKeyer {
	return norm.NewSinglePrimaryKey(field.Name("Id"))
}

// IsNew returns true if the model is new and wasn't loaded from storage
// Id must not be scanned for new models. Scanning in an Id effectively makes it appear not not new.
func (u *User) IsNew() bool {
	return !u.Id.Valid
}

var _ norm.Model = &User{}

// outputJson jsonEncode OR report error for use in Example Output: block
func outputJson(o interface{}) {
	str, err := json.Marshal(o)
	if err != nil {
		fmt.Println(err.Error())
	}
	os.Stdout.Write(str)
}

func ExampleModel() {

	user := &User{}
	user.Id.Scan("1")
	user.FirstName.Scan("John")
	user.LastName.Scan("Smith")
	user.Email.Scan("jsmith@example.org")

	str, _ := json.Marshal(user)

	os.Stdout.Write(str)
	// Output:
	// {"id":1,"first_name":"John","last_name":"Smith","email":"jsmith@example.org"}

}
