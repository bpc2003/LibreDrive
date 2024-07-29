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
	"github.com/golang-jwt/jwt/v5"
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
		enc.Encode(errStruct{Success: false, Msg: "No User with ID " + strconv.Itoa(userId)})
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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	id, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(errStruct{Success: false, Msg: "Invalid ID"})
		return
	}
	err = q.DeleteUser(ctx, int64(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		enc.Encode(errStruct{Success: false, Msg: "No user with ID " + strconv.Itoa(id)})
	} else {
		w.WriteHeader(http.StatusNoContent)
		enc.Encode(errStruct{Success: true, Msg: ""})
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

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginParams users.GetUserParams
	enc := json.NewEncoder(w)

	err := json.NewDecoder(r.Body).Decode(&loginParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(errStruct{Success: false, Msg: "Internal Error"})
		return
	}

	h := sha1.New()
	h.Write([]byte(loginParams.Password))
	loginParams.Password = base64.URLEncoding.EncodeToString(h.Sum(nil))
	
	user, err := q.GetUser(ctx, loginParams)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		enc.Encode(errStruct{Success: false, Msg: "Incorrect username or password"})
	} else {
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"Id": user.ID,
			"isAdmin": user.Isadmin,
		})
		tokString, _ := tok.SignedString([]byte(user.Password))
		w.Write([]byte(tokString))
	}
}
