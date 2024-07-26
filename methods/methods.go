package methods

import (
	"context"
	"crypto/sha1"
	"database/sql"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	enc := json.NewEncoder(w)
	users, err := q.GetUsers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(errStruct{Success: false, Msg: "Internal Error"})
	} else {
		enc.Encode(users)
	}
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(errStruct{Success: false, Msg: "Invalid ID"})
		return
	}

	user, err := q.GetUserById(ctx, int64(userId))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		enc.Encode(errStruct{Success: false, Msg: "No User with ID of " + strconv.Itoa(userId)})
	} else {
		enc.Encode(user)
	}
}

func ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	var passwordParams users.ChangePasswordParams

	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(errStruct{Success: false, Msg: "Invalid ID"})
		return
	}
	passwordParams.ID = int64(userId)

	err = json.NewDecoder(r.Body).Decode(&passwordParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(errStruct{Success: false, Msg: "Internal Error"})
	} else {
		h := sha1.New()
		h.Write([]byte(passwordParams.Password))
		passwordParams.Password = base64.URLEncoding.EncodeToString(h.Sum(nil))

		user, err := q.ChangePassword(ctx, passwordParams)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(errStruct{Success: false, Msg: "Internal Error"})
		} else {
			enc.Encode(user)
		}
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userParams users.CreateUserParams
	enc := json.NewEncoder(w)

	err := json.NewDecoder(r.Body).Decode(&userParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(errStruct{Success: false, Msg: "Internal Error"})
	} else {
		h := sha1.New()
		h.Write([]byte(userParams.Password))
		userParams.Password = base64.URLEncoding.EncodeToString(h.Sum(nil))

		user, err := q.CreateUser(ctx, userParams)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(errStruct{Success: false, Msg: "Internal Error"})
		} else {
			enc.Encode(user)
		}
	}
}
