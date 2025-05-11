package dfs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Step mewakili satu langkah kombinasi: A + B => C

// WriteDataMultiplePaths menyimpan hasil pencarian multiple recipe ke file JSON
func WriteDataMultiplePaths(
	pathResult string,
	elementName string,
	searchMethod string,
	isFound bool,
	allPaths []RecipePath,
	nodeCountElement int,
	searchTimeMS int,
) {
	// Siapkan struktur data hasil
	result := map[string]interface{}{
		"Element":                  elementName,
		"Method":                   searchMethod,
		"isFound":                  isFound,
		"NodeCountElement":         nodeCountElement,
		"SearchTimeInMilliseconds": searchTimeMS,
		"TotalPathsFound":          len(allPaths),
		"Paths":                    allPaths,
	}

	// Encode dan simpan ke file
	file, err := os.Create(pathResult)
	if err != nil {
		log.Fatal("Gagal membuat file output:", err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal("Gagal format JSON:", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatal("Gagal menulis file JSON:", err)
	}

	fmt.Printf("Berhasil menyimpan %d jalur pembentukan [%s] ke file [%s]\n", len(allPaths), elementName, pathResult)
}
