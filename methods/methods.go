package methods

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"libredrive/users"
)

var db *sql.DB
var q *users.Queries
var ctx = context.Background()

//go:embed schema.sql
var ddl string

func init() {
	var err error
	db, err = sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
	}
	q = users.New(db)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := q.GetUsers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Error"))
	} else {
		fmt.Fprintf(w, "%v\n", users)
	}
}
