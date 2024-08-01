package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"libredrive/models"
	"libredrive/types"
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

	passwordParams := models.ChangePasswordParams{}
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
