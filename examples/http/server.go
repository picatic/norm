package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/picatic/norm"
	"log"
	"net/http"
)

var dbrConnection *dbr.Connection

func main() {
	mysqlUri := flag.String("mysql", "norm_demo:password@(localhost:3306)/norm_demo?parseTime=true", "provide a database/sql driver uri")

	// configure sql driver and dbr
	db, err := sql.Open("mysql", *mysqlUri)
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
	http.HandleFunc("/users/[0-9]+", usersIdHandler)
	http.ListenAndServe(":8080", nil)
}

func usersHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("Hello"))
	switch r.Method {
	case "POST":
		norm.NewSelect(dbrConnection.NewSession(nil), &User{}, nil)
	case "GET":
		users := make([]*User, 0)

		selectQ := norm.NewSelect(dbrConnection.NewSession(nil), &User{}, nil)
		_, err := selectQ.LoadStructs(users)
		if err != nil {
			handleError(rw, err, http.StatusInternalServerError)
			return
		}
		out, err := json.Marshal(users)
		if err != nil {
			handleError(rw, err, http.StatusInternalServerError)
			return
		}
		rw.Write(out)
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
	default:
		rw.Write([]byte("WHAT?"))
	}
}

func handleError(rw http.ResponseWriter, err error, statusCode int) {
	rw.WriteHeader(statusCode)
	rw.Write([]byte(err.Error()))
}
