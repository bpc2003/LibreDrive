package controllers

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
	"libredrive/crypto"
	"libredrive/global"
	"libredrive/models"
)

var db *sql.DB
var q *models.Queries
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
	q = models.New(db)

	if u, _ := q.GetUsers(ctx); len(u) == 0 {
		password, salt := crypto.GeneratePassword(global.ADMIN_PASSWORD, 144)
		_, err = q.CreateUser(ctx, models.CreateUserParams{Username: "admin", Password: string(password), Salt: salt, Isadmin: true})
		if err != nil || os.MkdirAll(path.Join("users", "1"), 0750) != nil {
			log.Fatal(err)
		}
	}
}
