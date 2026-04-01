package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

// ========================================
// CHECKLIST 1: Batasi Penamaan File
// ========================================
// 1. Nama file tidak boleh kosong
// 2. Panjang nama file maksimal 255 karakter
// 3. Tidak boleh mengandung path traversal (.. / \)
// 4. Tidak boleh dimulai dengan titik (hidden files)
// 5. Hanya karakter aman: alfanumerik, -, _, ., dan spasi
// ========================================

var safeFileNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_\-. ]*$`)

func validateFileName(filename string) error {
	// 1. Nama file tidak boleh kosong
	if strings.TrimSpace(filename) == "" {
		return fmt.Errorf("nama file tidak boleh kosong")
	}

	// 2. Panjang nama file maksimal 255 karakter
	if len(filename) > 255 {
		return fmt.Errorf("nama file terlalu panjang (maksimal 255 karakter)")
	}

	// 3. Tidak boleh mengandung path traversal
	if strings.Contains(filename, "..") ||
		strings.Contains(filename, "/") ||
		strings.Contains(filename, "\\") {
		return fmt.Errorf("nama file tidak boleh mengandung path traversal (.. / \\)")
	}

	// 4. Tidak boleh dimulai dengan titik (hidden file)
	if strings.HasPrefix(filename, ".") {
		return fmt.Errorf("nama file tidak boleh dimulai dengan titik")
	}

	// 5. Hanya karakter aman yang diizinkan
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	if !safeFileNameRegex.MatchString(nameWithoutExt) {
		return fmt.Errorf("nama file hanya boleh mengandung huruf, angka, -, _, titik, dan spasi")
	}

	return nil
}

// ========================================
// CHECKLIST 2: Batasi Tipe File
// ========================================
// 1. Cek ekstensi file (whitelist)
// 2. Cek MIME type dari konten file (Content-Type)
// ========================================

// Daftar ekstensi yang diizinkan
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".txt":  true,
}

// Daftar MIME type yang diizinkan
var allowedMIMETypes = map[string]bool{
	"image/jpeg":         true,
	"image/png":          true,
	"image/gif":          true,
	"application/pdf":    true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"text/plain": true,
}

func validateFileType(filename string, fileHeader []byte, r *http.Request) error {
	// 1. Cek ekstensi file
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return fmt.Errorf("file harus memiliki ekstensi")
	}
	if !allowedExtensions[ext] {
		allowed := make([]string, 0, len(allowedExtensions))
		for k := range allowedExtensions {
			allowed = append(allowed, k)
		}
		return fmt.Errorf("ekstensi %s tidak diizinkan. Ekstensi yang diizinkan: %s", ext, strings.Join(allowed, ", "))
	}

	// 2. Cek MIME type dari konten file
	detectedType := http.DetectContentType(fileHeader)
	if !allowedMIMETypes[detectedType] {
		return fmt.Errorf("tipe file %s tidak diizinkan", detectedType)
	}

	return nil
}
