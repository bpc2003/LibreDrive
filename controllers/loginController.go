package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"libredrive/models"
	"libredrive/types"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	userParams := models.CreateUserParams{}
	r.ParseForm()
	userParams.Username = r.Form.Get("Username")
	isAdmin, err := strconv.ParseBool(r.Form.Get("IsAdmin"))
	if userParams.Username == "" || r.Form.Get("Password") == "" || err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userParams.Isadmin = isAdmin
	password, _ := bcrypt.GenerateFromPassword([]byte(r.Form.Get("Password")), 14)
	userParams.Password = string(password)

	if user, err := types.Queries.CreateUser(types.CTX, userParams); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	} else {
		os.MkdirAll(fmt.Sprintf("users/%d", user.ID), 0750)
		w.Write([]byte("Successfully created user"))
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Username := r.Form.Get("Username")
	Password := r.Form.Get("Password")

	user, err := types.Queries.GetUser(types.CTX, Username)
	if err != nil ||
		bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Password)) != nil {
		http.Error(w, "Incorrect Username or Password", http.StatusForbidden)
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
