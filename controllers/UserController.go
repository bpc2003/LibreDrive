package controllers

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/go-chi/chi/v5"
	"libredrive/crypto"
	"libredrive/global"
	"libredrive/models"
	"libredrive/templates"
)

// GetUsers - allows an admin user to see all the registered users.
func GetUsers(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAdmin").(bool) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	users, err := q.GetUsers(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		templates.Users(users).Render(ctx, w)
	}
}

func MarkUserActive(w http.ResponseWriter, r *http.Request) {
	var userId int64
	
	id := chi.URLParam(r, "id")
	if userId = global.ActiveTab[id]; userId == 0 {
		fmt.Println(global.ActiveTab)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	q.MarkActive(ctx, userId)
	delete(global.ActiveTab, id)
	w.Write([]byte("Account Activated"))
}

// ChangeUserPassword - allows a user to change their password.
func ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	passwordParams := models.ChangePasswordParams{}
	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	r.ParseForm()
	if err != nil || r.Form.Get("Password") == "" || len(r.Form.Get("Password")) > 72 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if r.Context().Value("id").(int) != userId {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	password, salt := crypto.GeneratePassword(r.Form.Get("Password"), 144)
	passwordParams.Password = password
	passwordParams.Salt = salt
	passwordParams.ID = int64(userId)

	if _, err := q.ChangePassword(ctx, passwordParams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		key := r.Context().Value("key").(string)
		nk := r.Form.Get("Password")
		var h [sha256.Size]byte
		for _, r := range salt {
			h = sha256.Sum256([]byte(string(r) + nk))
			nk = string(h[:])
		}
		h = sha256.Sum256([]byte(nk + salt))
		nk = fmt.Sprintf("%x", h)
		files, _ := os.ReadDir(path.Join("users", strconv.Itoa(userId)))
		for _, file := range files {
			f, _ := os.OpenFile(path.Join("users", strconv.Itoa(userId), file.Name()), os.O_RDWR, 0640)
			buf, _ := io.ReadAll(f)
			plain, _ := crypto.Decrypt([]byte(key), buf)
			cipher := crypto.Encrypt([]byte(nk), plain)
			f.Seek(0, 0)
			f.Write(cipher)
			f.Close()

			c := http.Cookie{
				Name:   "auth",
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			}
			http.SetCookie(w, &c)
			w.Header().Set("HX-Refresh", "true")
		}
	}
}

// ResetUserPassword - allows an admin user to reset another user's password.
// NOTE: When ResetUserPassword is called all previously encrypted files are lost.
func ResetUserPassword(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAdmin").(bool) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	passwordParams := models.ChangePasswordParams{}
	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	r.ParseForm()
	if err != nil || r.Form.Get("Password") == "" || len(r.Form.Get("Password")) > 72 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	password, salt := crypto.GeneratePassword(r.Form.Get("Password"), 14)
	passwordParams.Password = password
	passwordParams.Salt = salt
	passwordParams.ID = int64(userId)
	if _, err := q.ChangePassword(ctx, passwordParams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		files, _ := os.ReadDir(path.Join("users", strconv.Itoa(userId)))
		lost, _ := os.Create(path.Join("users", strconv.Itoa(userId), "lost.zip"))
		zw := zip.NewWriter(lost)
		for _, file := range files {
			f, _ := os.Open(path.Join("users", strconv.Itoa(userId), file.Name()))
			w, _ := zw.Create(file.Name())
			if _, err := io.Copy(w, f); err != nil {
				log.Fatal(err)
			}
			f.Close()
			os.Remove(path.Join("users", strconv.Itoa(userId), file.Name()))
		}
		zw.Close()
		w.Header().Set("HX-Refresh", "true")
	}
}

// DeleteUser - allows an admin to delete a user.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = q.DeleteUser(ctx, int64(userId))
	if err != nil {
		http.Error(w, fmt.Sprintf("No User with ID of %d", userId), http.StatusNotFound)
	} else {
		os.RemoveAll(path.Join("users", strconv.Itoa(userId)))
		w.Header().Set("HX-Refresh", "true")
	}
}
