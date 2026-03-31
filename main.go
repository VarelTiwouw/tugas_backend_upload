package main

import (
	"net/http"
	"tugas31/routes"
)

func main() {
	routes.RegisterRoutes()

	println("Server jalan di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
