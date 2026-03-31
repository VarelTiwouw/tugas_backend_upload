package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method harus POST", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println("ERROR:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Gagal ambil file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// buat folder uploads
	os.MkdirAll("uploads", os.ModePerm)

	filePath := filepath.Join("uploads", handler.Filename)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Gagal simpan file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	io.Copy(dst, file)

	w.Write([]byte("Upload berhasil"))
}
