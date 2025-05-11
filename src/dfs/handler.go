package dfs

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

/*
element = nama element yang ingin dicari resepnya
multiple = true jika ingin mencari multiple recipe dari sebuah elemen, vice versa
n = banyaknya resep berbeda yang ingin dicari, jika isMultiple false maka isi 1

Example: http://localhost:8080/dfs?element=Gold&multiple=true&n=5
*/

func DFSHandler(w http.ResponseWriter, r *http.Request) {
	// Ambil query parameter dari URL
	element := r.URL.Query().Get("element")
	multiple := r.URL.Query().Get("multiple") == "true"
	nStr := r.URL.Query().Get("n")

	// Validasi dan parsing jumlah recipe
	n := 1
	if multiple && nStr != "" {
		parsed, err := strconv.Atoi(nStr)
		if err == nil && parsed > 0 {
			n = parsed
		}
	}

	// Buat path file hasil (i.e. "./data/result.json")
	pathResult := filepath.Join("data", "result.json")

	// Jalankan DFS
	if multiple {
		DFS(element, pathResult, true, n)
	} else {
		DFS(element, pathResult, false, 0)
	}

	// Baca file JSON yang dihasilkan
	jsonBytes, err := os.ReadFile(pathResult)
	if err != nil {
		http.Error(w, "Gagal membaca file hasil", http.StatusInternalServerError)
		return
	}

	// Set header dan kirim response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}
