package types

import (
	"context"
  "database/sql"
  _ "embed"
  "log"

	_ "github.com/mattn/go-sqlite3"
  "libredrive/users"
)

var db *sql.DB
var Queries *users.Queries
var CTX = context.Background()

//go:embed schema.sql
var DDL string

func init() {
	var err error
	db, err = sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.ExecContext(CTX, DDL); err != nil {
		log.Fatal(err)
	}
	Queries = users.New(db)
}
