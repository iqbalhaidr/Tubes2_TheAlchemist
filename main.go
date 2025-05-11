package main

import (
	"littlealchemy/dfs" // Sesuaikan dengan import handler Anda
	"log"
	"net/http"
)

func main() {
	// Menyediakan handler untuk route /dfs
	http.HandleFunc("/dfs", dfs.DFSHandler)

	// Jalankan server di port 8080
	log.Println("Server berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
