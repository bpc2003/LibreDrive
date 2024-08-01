package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"libredrive/types"
	"github.com/go-chi/chi/v5"
)

func GetFiles(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	id := int(r.Context().Value("id").(float64))

	if files, err := os.ReadDir("users/" + strconv.Itoa(id)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(types.ErrStruct{Success: false, Msg: "Internal Error"})
	} else {
		fileNames := make([]string, 0)
		for _, f := range files {
			fileNames = append(fileNames, f.Name())
		}
		enc.Encode(fileNames)
	}
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	id := int(r.Context().Value("id").(float64))
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(types.ErrStruct{Success: false, Msg: err.Error()})
		return
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(types.ErrStruct{Success: false, Msg: err.Error()})
		return
	}

	if err = os.WriteFile("users/"+strconv.Itoa(id)+"/"+handler.Filename, buf, 0750); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(types.ErrStruct{Success: false, Msg: err.Error()})
	} else {
		w.WriteHeader(http.StatusNoContent)
		enc.Encode(types.ErrStruct{Success: true, Msg: ""})
	}
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	id := int(r.Context().Value("id").(float64))

	http.ServeFile(w, r, "users/"+strconv.Itoa(id)+"/"+fileName)
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	fileName := chi.URLParam(r, "fileName")
	id := int(r.Context().Value("id").(float64))

	if err := os.Remove("users/"+strconv.Itoa(id)+"/"+fileName); err != nil {
		w.WriteHeader(http.StatusNotFound)
		enc.Encode(types.ErrStruct{Success: false, Msg: "No file named " + fileName})
	} else {
		w.WriteHeader(http.StatusNoContent)
		enc.Encode(types.ErrStruct{Success: true, Msg: ""})
	}
}
