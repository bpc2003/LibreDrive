package controllers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/secretbox"
	"libredrive/types"
	"libredrive/templates"
)

func GetFiles(w http.ResponseWriter, r *http.Request) {
	id := int(r.Context().Value("id").(float64))

	if files, err := os.ReadDir(fmt.Sprintf("users/%d", id)); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	} else {
		fileNames := make([]string, 0)
		for _, f := range files {
			fileNames = append(fileNames, f.Name()[:len(f.Name())-4])
		}
		templates.Files(fileNames).Render(types.CTX, w)
	}
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	id := int(r.Context().Value("id").(float64))
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	user, _ := types.Queries.GetUserById(types.CTX, int64(id))
	key, _ := nacl.Load(fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password))))
	buf, _ := io.ReadAll(file)

	encrypted := secretbox.EasySeal(buf, key)
	if err = os.WriteFile(fmt.Sprintf("users/%d/%s.enc", id, handler.Filename), encrypted, 0750); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	id := int(r.Context().Value("id").(float64))
	user, _ := types.Queries.GetUserById(types.CTX, int64(id))
	key, _ := nacl.Load(fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password))))

	fp, err := os.Open(fmt.Sprintf("users/%d/%s.enc", id, fileName))
	if err != nil {
		http.Error(w, fmt.Sprintf("File '%s' doesn't exist", fileName), http.StatusNotFound)
		return
	}
	defer fp.Close()
	buf, _ := io.ReadAll(fp)

	if buf, err = secretbox.EasyOpen(buf, key); err != nil {
		log.Fatal(err)
	} else {
		fp, _ = os.Create(fmt.Sprintf("users/%d/%s", id, fileName))
		defer fp.Close()
		io.WriteString(fp, string(buf))
		http.ServeFile(w, r, fmt.Sprintf("users/%d/%s", id, fileName))
		os.Remove(fmt.Sprintf("users/%d/%s", id, fileName))
	}
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	id := int(r.Context().Value("id").(float64))

	if err := os.Remove(fmt.Sprintf("users/%d/%s.enc", id, fileName)); err != nil {
		http.Error(w, fmt.Sprintf("File '%s' doesn't exist", fileName), http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
