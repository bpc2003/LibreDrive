package controllers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/secretbox"
	"golang.org/x/crypto/bcrypt"
	"libredrive/models"
	"libredrive/templates"
	"libredrive/types"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := types.Queries.GetUsers(types.CTX)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		templates.Users(users).Render(types.CTX, w)
	}
}

func ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	passwordParams := models.ChangePasswordParams{}
	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	r.ParseForm()
	if err != nil || r.Form.Get("Password") == "" || len(r.Form.Get("Password")) > 72 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(r.Form.Get("Password")), 14)
	passwordParams.Password = string(password)
	passwordParams.ID = int64(userId)

	if _, err := types.Queries.ChangePassword(types.CTX, passwordParams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		key, _ := nacl.Load(r.Context().Value("key").(string))
		nk, _ := nacl.Load(fmt.Sprintf("%x", sha256.Sum256([]byte(passwordParams.Password))))
		files, _ := os.ReadDir(path.Join("user_data", strconv.Itoa(userId)))
		for _, file := range files {
			f, _ := os.OpenFile(path.Join("user_data", strconv.Itoa(userId), file.Name()), os.O_RDWR, 0640)
			buf, _ := io.ReadAll(f)
			plain, _ := secretbox.EasyOpen(buf, key)
			cipher := secretbox.EasySeal(plain, nk)
			f.Seek(0, 0)
			f.Write(cipher)
			f.Close()
		}
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = types.Queries.DeleteUser(types.CTX, int64(userId))
	if err != nil {
		http.Error(w, fmt.Sprintf("No User with ID of %d", userId), http.StatusNotFound)
	} else {
		os.RemoveAll(path.Join("user_data", strconv.Itoa(userId)))
		w.Header().Set("HX-Refresh", "true")
	}
}
