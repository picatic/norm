package norm_test

import (
	"fmt"
	"github.com/picatic/norm"
	"github.com/picatic/norm/field"
	"github.com/gocraft/dbr"
	"github.com/gocraft/dbr/dialect"
)

func ExampleModelFields() {

	user := &User{}
	fmt.Println(norm.ModelFields(user))
	// Output:
	// [Id FirstName LastName Email]

}

func ExampleModelGetField_() {

	user := &User{}
	user.Id.Scan(1234)
	modelField, err := norm.ModelGetField(user, field.Name("Id"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	outputJson(modelField)

	// Output:
	// 1234
}

func ExampleNewSelect() {

	cnx := norm.NewConnection(nil, "norm_mysql", nil)
	dbrSess := cnx.NewSession(nil)
	buf := dbr.NewBuffer()
	user := &User{}
	user.Id.Scan(1)

	selectBuilder := norm.NewSelect(dbrSess, user, nil)
	err := selectBuilder.Build(dialect.MySQL, buf)
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("DML:", buf)

	// Output:
	// DML: SELECT `id`, `first_name`, `last_name`, `email` FROM norm_mysql.norm_users

}

func SkipExampleNewUpdate() {

	cnx := norm.NewConnection(nil, "norm_mysql", nil)
	dbrSess := cnx.NewSession(nil)
	buf := dbr.NewBuffer()

	user := &User{}
	user.Id.Scan(1)
	user.FirstName.Scan("Zim")
	user.LastName.Scan("Ham")
	user.Email.Scan("zh@example.com")

	updateBuilder := norm.NewUpdate(dbrSess, user, nil).Where("id = ?", user.Id.Int64)
	err := updateBuilder.Build(dialect.MySQL, buf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("DML:", buf.String())
	fmt.Print("ARGS: ")
	outputJson(buf.Value())

	// https://github.com/picatic/norm/issues/4
	// Unstable Field Order

	// Output:
	// DML: UPDATE norm_mysql.norm_users SET `first_name` = ?, `last_name` = ?, `email` = ? WHERE (id = ?)
	// ARGS: ["Zim","Ham","zh@example.com",1]

}

func ExampleNewInsert() {

	cnx := norm.NewConnection(nil, "norm_mysql", nil)
	dbrSess := cnx.NewSession(nil)
	buf := dbr.NewBuffer()
	user := &User{}
	user.FirstName.Scan("Zim")
	user.LastName.Scan("Ham")
	user.Email.Scan("zh@example.com")

	insertBuilder := norm.NewInsert(dbrSess, user, nil).Record(user)
	err := insertBuilder.Build(dialect.MySQL, buf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("DML:", buf.String())
	fmt.Print("ARGS: ")
	outputJson(buf.Value())

	// Output:
	// DML: INSERT INTO `norm_mysql`.`norm_users` (`first_name`,`last_name`,`email`) VALUES (?,?,?)
	// ARGS: ["Zim","Ham","zh@example.com"]

}

func ExampleNewDelete() {

	cnx := norm.NewConnection(nil, "norm_mysql", nil)
	dbrSess := cnx.NewSession(nil)
	buf := dbr.NewBuffer()
	user := &User{}
	user.Id.Scan(5432)

	deleteBuilder := norm.NewDelete(dbrSess, user).Where("id = ?", user.Id.Int64)
	err := deleteBuilder.Build(dialect.MySQL, buf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("DML:", buf.String())
	fmt.Print("ARGS: ")
	outputJson(buf.Value())

	// Output:
	// DML: DELETE FROM `norm_mysql`.`norm_users` WHERE (id = ?)
	// ARGS: [5432]

}

func ExampleModelLoadMap() {

	dataMap := map[string]interface{}{
		"id":         "1234",
		"first_name": "James",
	}
	user := &User{}

	norm.ModelLoadMap(user, dataMap)

	outputJson(user)
	// Output:
	// {"id":1234,"first_name":"James","last_name":"","email":""}
}

func ExampleModelDirtyFields() {

	user := &User{}

	user.FirstName.Scan("Roberto")

	dirtyFields, err := norm.ModelDirtyFields(user)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	fmt.Println("CLEAN:", dirtyFields)

	user.FirstName.Scan("Robbie")

	dirtyFields, err = norm.ModelDirtyFields(user)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	fmt.Println("DIRTY:", dirtyFields)

	// Output:
	// CLEAN: []
	// DIRTY: [FirstName]

}

func ExampleModel_TableName() {
	user := &User{}

	fmt.Println(user.TableName())
	// Output:
	// norm_users

}
