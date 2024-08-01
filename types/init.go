package types

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"libredrive/models"
)

var db *sql.DB
var Queries *models.Queries
var CTX = context.Background()

//go:embed schema.sql
var ddl string

func init() {
	var err error
	db, err = sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.ExecContext(CTX, ddl); err != nil {
		log.Fatal(err)
	}
	Queries = models.New(db)

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	if items, _ := Queries.GetUsers(CTX); len(items) == 0 {
		password, _ := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), 14)
		_, err := Queries.CreateUser(CTX, models.CreateUserParams{Username: "admin", Password: string(password), Isadmin: true})
		if err != nil || os.MkdirAll("users/1", 0750) != nil {
			log.Fatal(err)
		}
	}
}
