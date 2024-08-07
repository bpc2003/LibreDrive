package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	passwordParams := models.ChangePasswordParams{}
	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	r.ParseForm()
	if r.Form.Get("Password") == "" || len(r.Form.Get("Password")) > 72 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	passwordParams.ID = int64(userId)
	password, _ := bcrypt.GenerateFromPassword([]byte(r.Form.Get("Password")), 14)
	passwordParams.Password = string(password)
	user, err := types.Queries.GetUserById(types.CTX, passwordParams.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("No user with ID of %d", passwordParams.ID), http.StatusNotFound)
		return
	}

	if nu, err := types.Queries.ChangePassword(types.CTX, passwordParams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		key, _ := nacl.Load(fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password))))
		nk, _ := nacl.Load(fmt.Sprintf("%x", sha256.Sum256([]byte(nu.Password))))
		files, _ := os.ReadDir(fmt.Sprintf("users/%d", userId))
		for _, file := range files {
			f, _ := os.OpenFile(fmt.Sprintf("users/%d/%s", userId, file.Name()), os.O_RDWR, 0750)
			defer f.Close()

			buf, _ := io.ReadAll(f)
			plain, _ := secretbox.EasyOpen(buf, key)
			cipher := secretbox.EasySeal(plain, nk)
			f.Seek(0, 0)
			f.Write(cipher)
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
		w.WriteHeader(http.StatusNoContent)
		os.RemoveAll(fmt.Sprintf("users/%d", userId))
	}
}
