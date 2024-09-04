package controllers

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path"
	"strconv"
	"strings"

	"libredrive/crypto"
	"libredrive/global"
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
	userParams.Email = r.Form.Get("Email")
	userParams.Isadmin = r.Form.Get("IsAdmin") == "on"
	userParams.Active = false
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
		global.ActiveTab[user.Password] = user.ID
		err := smtp.SendMail(global.AUTH_HOST+":"+global.AUTH_PORT,
			global.Auth, global.AUTH_EMAIL,
			[]string{user.Email},
			[]byte("Activate Your Account Here: http://" + global.HOST + ":" + global.PORT + "/api/activate/" + user.Password))
		if err != nil {
			log.Fatal(err.Error())
		}
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
		crypto.ComparePassword(Password, user.Salt, user.Password) == false ||
		user.Active == false {
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
