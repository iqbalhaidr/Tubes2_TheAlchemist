package main

import (
	"encoding/json"
	"littlealchemy/src/bfs"
	"littlealchemy/src/dfs"
	"log"
	"net/http"
)

type Message struct {
	Message string `json:"message"`
}

// untuk cek hubungan ke Golang dari React

// functions for ease
func withCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		withCORS(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	withCORS(w)
	w.Header().Set("Content-Type", "application/json")
	message := Message{Message: "Hello World! From Golang!"}
	json.NewEncoder(w).Encode(message)
}

func dfsWithCORS(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		withCORS(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	withCORS(w)
	dfs.DFSHandler(w, r)
}

func bfsWithCORS(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		withCORS(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	withCORS(w)
	bfs.BFSHandler(w, r)
}

func main() {
	/*
		element = nama element yang ingin dicari resepnya
		multiple = true jika ingin mencari multiple recipe dari sebuah elemen, vice versa
		n = banyaknya resep berbeda yang ingin dicari, jika isMultiple false maka isi 1

		Example: http://localhost:8080/dfs?element=Gold&multiple=true&n=5
	*/

	router := http.NewServeMux()

	// handler untuk cek hubungan ke Golang dari React
	router.HandleFunc("/ping", handler)

	// Menyediakan handler untuk route /dfs
	router.HandleFunc("/dfs", dfsWithCORS)

	// Menyediakan handler untuk route /bfs
	router.HandleFunc("/bfs", bfsWithCORS)

	// Jalankan server di port 8080
	log.Println("Server berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
