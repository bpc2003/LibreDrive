package controllers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/secretbox"
	"libredrive/templates"
	"libredrive/types"
)

func GetFiles(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	if files, err := os.ReadDir(path.Join("users", strconv.Itoa(id))); err != nil {
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
	id := r.Context().Value("id").(int)
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
	if err = os.WriteFile(path.Join("users", strconv.Itoa(id), handler.Filename+".enc"), encrypted, 0750); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("HX-Refresh", "true")
	}
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	id := r.Context().Value("id").(int)
	user, _ := types.Queries.GetUserById(types.CTX, int64(id))
	key, _ := nacl.Load(fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password))))

	fp, err := os.Open(path.Join("users", strconv.Itoa(id), fileName+".enc"))
	if err != nil {
		http.Error(w, fmt.Sprintf("File '%s' doesn't exist", fileName), http.StatusNotFound)
		return
	}
	defer fp.Close()
	buf, _ := io.ReadAll(fp)

	if buf, err = secretbox.EasyOpen(buf, key); err != nil {
		log.Fatal(err)
	} else {
		fp, _ = os.Create(path.Join("users", strconv.Itoa(id), fileName))
		fp.Write(buf)
		fp.Close()

		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, path.Join("users", strconv.Itoa(id), fileName))
		os.Remove(path.Join("users", strconv.Itoa(id), fileName))
	}
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	id := r.Context().Value("id").(int)

	if err := os.Remove(path.Join("users", strconv.Itoa(id), fileName+".enc")); err != nil {
		http.Error(w, fmt.Sprintf("File '%s' doesn't exist", fileName), http.StatusNotFound)
	} else {
		w.Header().Set("HX-Refresh", "true")
	}
}
