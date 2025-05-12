package bfs

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func BFSHandler(w http.ResponseWriter, r *http.Request) {
	element := r.URL.Query().Get("element")
	multiple := r.URL.Query().Get("multiple") == "true"
	nStr := r.URL.Query().Get("n")

	n := 1
	if multiple && nStr != "" {
		parsed, err := strconv.Atoi(nStr)
		if err == nil && parsed > 0 {
			n = parsed
		}
	}

	pathResult := filepath.Join("data", "result.json")

	if multiple {
		BFS(element, pathResult, true, n)
	} else {
		BFS(element, pathResult, false, 0)
	}

	jsonBytes, err := os.ReadFile(pathResult)
	if err != nil {
		http.Error(w, "Gagal membaca file hasil", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}
