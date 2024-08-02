package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/secretbox"
	"libredrive/types"
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

	file, handler, err := r.FormFile("upload")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(types.ErrStruct{Success: false, Msg: err.Error()})
		return
	}
	defer file.Close()

	user, _ := types.Queries.GetUserById(types.CTX, int64(id))
	key, err := nacl.Load(fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password))))
	if err != nil {
		log.Fatal(err)
	}

	buf, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(types.ErrStruct{Success: false, Msg: err.Error()})
		return
	}

	encrypted := secretbox.EasySeal(buf, key)
	if err = os.WriteFile("users/"+strconv.Itoa(id)+"/"+handler.Filename+".enc", encrypted, 0750); err != nil {
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

	user, _ := types.Queries.GetUserById(types.CTX, int64(id))
	key, err := nacl.Load(fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password))))
	if err != nil {
		log.Fatal(err)
	}

	fp, err := os.Open("users/" + strconv.Itoa(id) + "/" + fileName + ".enc")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("File doesn't exist"))
		return
	}
	defer fp.Close()

	buf, err := io.ReadAll(fp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Error"))
		return
	}

	if buf, err = secretbox.EasyOpen(buf, key); err != nil {
		log.Fatal(err)
	} else {
		fp, _ = os.Create("users/" + strconv.Itoa(id) + "/" + fileName)
		defer fp.Close()
		io.WriteString(fp, string(buf))
		http.ServeFile(w, r, "users/"+strconv.Itoa(id)+"/"+fileName)
		os.Remove("users/" + strconv.Itoa(id) + "/" + fileName)
	}
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	fileName := chi.URLParam(r, "fileName")
	id := int(r.Context().Value("id").(float64))

	if err := os.Remove("users/" + strconv.Itoa(id) + "/" + fileName); err != nil {
		w.WriteHeader(http.StatusNotFound)
		enc.Encode(types.ErrStruct{Success: false, Msg: "No file named " + fileName})
	} else {
		w.WriteHeader(http.StatusNoContent)
		enc.Encode(types.ErrStruct{Success: true, Msg: ""})
	}
}
