package routes

import (
	"net/http"
	"tugas31/handlers"
)

func RegisterRoutes() {
	http.HandleFunc("/upload", handlers.UploadFile)
}
