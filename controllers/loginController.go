package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"libredrive/types"
	"libredrive/users"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	userParams := users.CreateUserParams{}
	enc := json.NewEncoder(w)

	err := json.NewDecoder(r.Body).Decode(&userParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(types.ErrStruct{Success: false, Msg: "Internal Error"})
	} else {
		password, _ := bcrypt.GenerateFromPassword([]byte(userParams.Password), 14)
		userParams.Password = string(password)

		user, err := types.Queries.CreateUser(types.CTX, userParams)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(types.ErrStruct{Success: false, Msg: "Internal Error"})
		} else {
			if err := os.MkdirAll("user_data/"+strconv.Itoa(int(user.ID)), 0750); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				enc.Encode(types.ErrStruct{Success: false, Msg: err.Error()})
				return
			}
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
