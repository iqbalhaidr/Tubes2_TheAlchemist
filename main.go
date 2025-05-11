package main

import (
	"littlealchemy/src/dfs" // Sesuaikan dengan import handler Anda
	"log"
	"net/http"
)

func main() {
	/*
		element = nama element yang ingin dicari resepnya
		multiple = true jika ingin mencari multiple recipe dari sebuah elemen, vice versa
		n = banyaknya resep berbeda yang ingin dicari, jika isMultiple false maka isi 1

		Example: http://localhost:8080/dfs?element=Gold&multiple=true&n=5
	*/

	// Menyediakan handler untuk route /dfs
	http.HandleFunc("/dfs", dfs.DFSHandler)

	// Jalankan server di port 8080
	log.Println("Server berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
