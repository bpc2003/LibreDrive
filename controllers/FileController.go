// controllers - controllers for all the various routes
package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/go-chi/chi/v5"
	"libredrive/crypto"
	"libredrive/templates"
)

// GetFiles - Gets a user's files.
func GetFiles(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	if files, err := os.ReadDir(path.Join("users", strconv.Itoa(id))); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	} else {
		fileNames := make([]string, 0)
		for _, f := range files {
			fileNames = append(fileNames, f.Name())
		}
		templates.Files(fileNames).Render(ctx, w)
	}
}

// UploadFile - takes a user file encrypts it, and then stores it with the rest of their files.
func UploadFile(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	key := r.Context().Value("key").(string)
	buf, _ := io.ReadAll(file)

	encrypted := crypto.Encrypt([]byte(key), buf)
	if err = os.WriteFile(path.Join("users", strconv.Itoa(id), handler.Filename), encrypted, 0640); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("HX-Refresh", "true")
	}
}

// GetFile - allows a user to download a file.
func GetFile(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	id := r.Context().Value("id").(int)
	key := r.Context().Value("key").(string)
	fp, err := os.Open(path.Join("users", strconv.Itoa(id), fileName))
	if err != nil {
		http.Error(w, fmt.Sprintf("File '%s' doesn't exist", fileName), http.StatusNotFound)
		return
	}
	defer fp.Close()
	buf, _ := io.ReadAll(fp)

	if fileName != "lost.zip" {
		if buf, err = crypto.Decrypt([]byte(key), buf); err != nil {
			log.Fatal(err)
		}
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(buf)
}

// DeleteFile - allows a user to delete a file.
func DeleteFile(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	id := r.Context().Value("id").(int)

	if err := os.Remove(path.Join("users", strconv.Itoa(id), fileName)); err != nil {
		http.Error(w, fmt.Sprintf("File '%s' doesn't exist", fileName), http.StatusNotFound)
	} else {
		w.Header().Set("HX-Refresh", "true")
	}
}
