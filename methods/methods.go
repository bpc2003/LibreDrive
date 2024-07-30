package methods

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"libredrive/types"
	"libredrive/users"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	users, err := types.Queries.GetUsers(types.CTX)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(types.ErrStruct{Success: false, Msg: "Internal Error"})
	} else {
		enc.Encode(users)
	}
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(types.ErrStruct{Success: false, Msg: "Invalid ID"})
		return
	}

	user, err := types.Queries.GetUserById(types.CTX, int64(userId))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		enc.Encode(types.ErrStruct{Success: false, Msg: "No User with ID " + strconv.Itoa(userId)})
	} else {
		enc.Encode(user)
	}
}

func ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(types.ErrStruct{Success: false, Msg: "Invalid ID"})
		return
	}

	passwordParams := users.ChangePasswordParams{}
	if err = json.NewDecoder(r.Body).Decode(&passwordParams); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(types.ErrStruct{Success: false, Msg: "Internal Error"})
		return
	}
	passwordParams.ID = int64(userId)
	password, _ := bcrypt.GenerateFromPassword([]byte(passwordParams.Password), 14)
	passwordParams.Password = string(password)

	user, err := types.Queries.ChangePassword(types.CTX, passwordParams)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		enc.Encode(types.ErrStruct{Success: false, Msg: "No user with ID " + strconv.Itoa(userId)})
	} else {
		enc.Encode(user)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(types.ErrStruct{Success: false, Msg: "Invalid ID"})
		return
	}

	err = types.Queries.DeleteUser(types.CTX, int64(userId))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		enc.Encode(types.ErrStruct{Success: false, Msg: "No user with ID " + strconv.Itoa(userId)})
	} else {
		w.WriteHeader(http.StatusNoContent)
		enc.Encode(types.ErrStruct{Success: true, Msg: ""})
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	userParams := users.CreateUserParams{}
	enc := json.NewEncoder(w)

	err := json.NewDecoder(r.Body).Decode(&userParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(types.ErrStruct{Success: false, Msg: "Internal Error"})
	} else {
		if err := os.MkdirAll("user_data/"+userParams.Username, 0750); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(types.ErrStruct{Success: false, Msg: err.Error()})
			return
		}
		password, _ := bcrypt.GenerateFromPassword([]byte(userParams.Password), 14)
		userParams.Password = string(password)

		user, err := types.Queries.CreateUser(types.CTX, userParams)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(types.ErrStruct{Success: false, Msg: "Internal Error"})
		} else {
			enc.Encode(user)
		}
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	loginParams := types.LoginParams{}
	if err := json.NewDecoder(r.Body).Decode(&loginParams); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(types.ErrStruct{Success: false, Msg: "Internal Error"})
		return
	}

	user, err := types.Queries.GetUser(types.CTX, loginParams.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginParams.Password)) != nil {
		w.WriteHeader(http.StatusForbidden)
		enc.Encode(types.ErrStruct{Success: false, Msg: "Incorrect username or password"})
		return
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     time.Now().Add(time.Duration(time.Minute * 30)).Unix(),
		"iat":     time.Now().Unix(),
		"id":      user.ID,
		"isAdmin": user.Isadmin,
	})
	tokString, _ := tok.SignedString([]byte(os.Getenv("SECRET")))
	w.Write([]byte(tokString))
}
