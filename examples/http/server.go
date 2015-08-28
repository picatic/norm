package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/picatic/norm"
	"log"
	"net/http"
	"net/url"
)

var dbrConnection *dbr.Connection

func main() {
	mysqlUri := flag.String("mysql", "mysql://norm_demo:password@localhost:3306/norm_demo", "provide a database/sql driver uri")
	flag.Parse()
	dbUri, _ := url.Parse(*mysqlUri)
	// configure sql driver and dbr
	var dbUriString = ""

	username := dbUri.User.Username()
	password, hasPassword := dbUri.User.Password()
	if hasPassword {
		dbUriString = fmt.Sprintf("%s:%s@", username, password)
	} else {
		dbUriString = fmt.Sprintf("%s@", username)
	}

	dbUriString = fmt.Sprintf("%s(%s)%s?parseTime=true", dbUriString, dbUri.Host, dbUri.Path)
	log.Printf("uri %s", dbUriString)
	db, err := sql.Open("mysql", dbUriString)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	dbrConnection = dbr.NewConnection(db, nil)

	// register some handlers
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/[0-9]?", usersIdHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func usersHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "POST":
		norm.NewSelect(dbrConnection.NewSession(nil), &User{}, nil)
	case "GET":
		users := make([]*User, 0)

		selectQ := norm.NewSelect(dbrConnection.NewSession(nil), &User{}, nil)
		_, err := selectQ.LoadStructs(&users)
		if err != nil {
			handleError(rw, err, http.StatusInternalServerError)
			return
		}
		handleJSON(rw, users)
	default:
		rw.Write([]byte("WHAT?"))
	}
	//norm.NewSelect(dbrConnection.NewSession(), User{}, nil).Where(Where("id = ?"))
}

func usersIdHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "POST":
		norm.NewSelect(dbrConnection.NewSession(nil), &User{}, nil)
	case "GET":
		user := &User{}
		norm.NewSelect(dbrConnection.NewSession(nil), &User{}, nil).Where("id = ?", 1).LoadStruct(user)

	default:
		rw.Write([]byte("WHAT?"))
	}
}

func handleJSON(rw http.ResponseWriter, data interface{}) {
	out, err := json.Marshal(data)
	if err != nil {
		handleError(rw, err, http.StatusInternalServerError)
		return
	}
	rw.Write(out)
}

func handleError(rw http.ResponseWriter, err error, statusCode int) {
	rw.WriteHeader(statusCode)
	rw.Write([]byte(err.Error()))
}
