package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"libredrive/models"
	"libredrive/types"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	userParams := models.CreateUserParams{}
	r.ParseForm()
	isAdmin, err := strconv.ParseBool(r.Form.Get("IsAdmin"))
	if r.Form.Get("Username") == "" || r.Form.Get("Password") == "" || len(r.Form.Get("Password")) > 72 || err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userParams.Username = r.Form.Get("Username")
	userParams.Isadmin = isAdmin
	password, _ := bcrypt.GenerateFromPassword([]byte(r.Form.Get("Password")), 14)
	userParams.Password = string(password)

	if user, err := types.Queries.CreateUser(types.CTX, userParams); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	} else {
		os.MkdirAll(path.Join("users", strconv.Itoa(int(user.ID))), 0750)
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
	c := http.Cookie{
		Name:   "auth",
		Value:  fmt.Sprintf("%d&%t", user.ID, user.Isadmin),
		MaxAge: 1800,
		Path:   "/",
	}
	http.SetCookie(w, &c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
