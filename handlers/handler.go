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

	// ============================================
	// CHECKLIST 1: Validasi penamaan file
	// ============================================
	if err := validateFileName(handler.Filename); err != nil {
		http.Error(w, fmt.Sprintf("Nama file tidak valid: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// ============================================
	// CHECKLIST 2: Validasi tipe file
	// ============================================
	// Baca 512 byte pertama untuk deteksi MIME type
	fileHeader := make([]byte, 512)
	n, err := file.Read(fileHeader)
	if err != nil {
		http.Error(w, "Gagal membaca file", http.StatusBadRequest)
		return
	}
	fileHeader = fileHeader[:n]

	// Reset posisi file ke awal setelah membaca header
	file.Seek(0, 0)

	if err := validateFileType(handler.Filename, fileHeader, r); err != nil {
		http.Error(w, fmt.Sprintf("Tipe file tidak valid: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Buat folder uploads
	os.MkdirAll("uploads", os.ModePerm)

	// Gunakan nama file yang sudah divalidasi
	safeFileName := filepath.Base(handler.Filename)
	filePath := filepath.Join("uploads", safeFileName)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Gagal simpan file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	io.Copy(dst, file)

	w.Write([]byte(fmt.Sprintf("Upload berhasil! File disimpan: %s", safeFileName)))
}
