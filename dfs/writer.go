package dfs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type ResultFormat struct {
	Element                 string      `json:"Element"`
	Method                  string      `json:"Method"`
	IsFound                 bool        `json:"isFound"`
	NodeCountElement        int         `json:"NodeCountElement"`
	NodeCountRecipe         int         `json:"NodeCountRecipe"`
	SearchTimeInMiliseconds int         `json:"SearchTimeInMiliseconds"`
	Steps                   [][3]string `json:"Steps"`
}

func WriteData(pathResult string, elementName string, searchMethod string, isFound bool, steps [][3]string, nodeCountElement int, nodeCountRecipe int, searchTimeMS int) {
	result := ResultFormat{
		Element:                 elementName,
		Method:                  searchMethod,
		IsFound:                 isFound,
		NodeCountElement:        nodeCountElement,
		NodeCountRecipe:         nodeCountRecipe,
		SearchTimeInMiliseconds: searchTimeMS,
		Steps:                   steps,
	}

	file, err := os.Create(pathResult)
	if err != nil {
		log.Fatal("Gagal membuat file output:", err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal("Gagal format JSON:", err)
	}

	file.Write(jsonData)

	fmt.Printf("Berhasil menyimpan hasil pencarian [%s] menggunakan metode [%s] di [%s]\n", elementName, searchMethod, pathResult)
}

/* ===================================================================================================================================================================== */

type ResultFormatN struct {
	Element                 string       `json:"Element"`
	Method                  string       `json:"Method"`
	IsFound                 bool         `json:"isSatisfied"`
	NodeCountElement        int          `json:"NodeCountElement"`
	NodeCountRecipe         int          `json:"NodeCountRecipe"`
	SearchTimeInMiliseconds int          `json:"SearchTimeInMiliseconds"`
	TotalPathWanted         int          `json:"TotalPathsWanted"`
	TotalPathsFound         int          `json:"TotalPathsFound"`
	Steps                   []RecipePath `json:"Steps"`
}

// WriteDataMultiplePaths menyimpan hasil pencarian multiple recipe ke file JSON
func WriteDataMultiplePaths(
	pathResult string,
	elementName string,
	searchMethod string,
	isFound bool,
	allPaths []RecipePath,
	nodeCountElement int,
	nodeCountRecipe int,
	searchTimeMS int,
	nPath int,
) {

	// Siapkan struktur data hasil
	result := ResultFormatN{
		Element:                 elementName,
		Method:                  searchMethod,
		IsFound:                 isFound,
		NodeCountElement:        nodeCountElement,
		NodeCountRecipe:         nodeCountRecipe,
		SearchTimeInMiliseconds: searchTimeMS,
		TotalPathWanted:         nPath,
		TotalPathsFound:         len(allPaths),
		Steps:                   allPaths,
	}

	file, err := os.Create(pathResult)
	if err != nil {
		log.Fatal("Gagal membuat file output:", err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(result, "", "  ")
	// jsonData, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Gagal format JSON:", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatal("Gagal menulis file JSON:", err)
	}

	fmt.Printf("Berhasil menyimpan [%d] jalur pembentukan [%s] ke file [%s]\n", len(allPaths), elementName, pathResult)
}
