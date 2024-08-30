package controllers

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"libredrive/crypto"
	"libredrive/models"
)

// CreateUser - allows an admin to create a user.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	userParams := models.CreateUserParams{}
	r.ParseForm()
	if r.Form.Get("Username") == "" || r.Form.Get("Password") == "" ||
		len(r.Form.Get("Password")) > 72 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userParams.Username = r.Form.Get("Username")
	userParams.Isadmin = r.Form.Get("IsAdmin") == "on"
	password, salt := crypto.GeneratePassword(r.Form.Get("Password"), 144)
	userParams.Password = password
	userParams.Salt = salt

	if userParams.Isadmin {
		auth, _ := r.Cookie("auth")
		if auth == nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		attrs := strings.Split(auth.Value, "&")
		isAdmin, _ := strconv.ParseBool(attrs[1])
		if !isAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	if user, err := q.CreateUser(ctx, userParams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		os.MkdirAll(path.Join("users", strconv.Itoa(int(user.ID))), 0750)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// LoginUser - allows a user to login to their account.
func LoginUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Username := r.Form.Get("Username")
	Password := r.Form.Get("Password")

	user, err := q.GetUser(ctx, Username)
	if err != nil ||
		crypto.ComparePassword(Password, user.Salt, user.Password) == false {
		http.Error(w, "Incorrect Username or Password", http.StatusForbidden)
		return
	}
	var h [sha256.Size]byte
	for _, r := range user.Salt {
		h = sha256.Sum256([]byte(string(r) + Password))
		Password = string(h[:])
	}
	h = sha256.Sum256([]byte(Password + user.Salt))
	c := http.Cookie{
		Name:   "auth",
		Value:  fmt.Sprintf("%d&%t&%x", user.ID, user.Isadmin, h),
		MaxAge: 1800,
		Path:   "/",
	}
	http.SetCookie(w, &c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
