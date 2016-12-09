package norm_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/picatic/norm"
	"github.com/picatic/norm/field"
	"log"
)

func ExampleConnection_selectBySql() {

	user := &User{}

	// Mock Database Connection to expect a query and return row data
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(fmt.Printf("an error'%s' was not expected while open mock database", err))
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow("5432")

	mock.ExpectQuery("SELECT 1").WillReturnRows(rows)

	// Create norm Connection and Session for sqlmock
	connection := norm.NewConnection(db, "norm_database", nil)
	session := connection.NewSession(nil)

	// Query and load mocked result into Model
	err = session.SelectBySql("SELECT 1").LoadStruct(user)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Compare query executed versus query provided
	err = mock.ExpectationsWereMet()
	if err != nil {
		fmt.Println(err.Error())
	}

	outputJson(user)
	// Output:
	// {"id":5432,"first_name":"","last_name":"","email":""}

}

func ExampleConnection_newSelect() {

	user := &User{}

	// Mock Database Connection to expect a query and return row data
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(fmt.Printf("an error'%s' was not expected while open mock database", err))
	}

	rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(5432, "soup@example.com")

	// ExpectQuery is a regexp, escape regex tokens like `\\(` and `\\)`
	mock.ExpectQuery("SELECT `id`, `email` FROM norm_database.norm_users WHERE \\(id = 5432\\)").WillReturnRows(rows)

	// Create norm Connection and Session for sqlmock
	connection := norm.NewConnection(db, "norm_database", nil)
	session := connection.NewSession(nil)

	// Perform and load mocked result into Model
	err = norm.NewSelect(session, user, field.Names{"id", "email"}).Where("id = ?", 5432).LoadStruct(user)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Compare query executed versus query provided
	err = mock.ExpectationsWereMet()
	if err != nil {
		fmt.Println(err.Error())
	}

	outputJson(user)
	// Output:
	// {"id":5432,"first_name":"","last_name":"","email":"soup@example.com"}

}
