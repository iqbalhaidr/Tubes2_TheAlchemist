package dfs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type SearchResult struct {
	Element                 string      `json:"Element"`
	Method                  string      `json:"Method"`
	IsFound                 bool        `json:"isFound"`
	NodeCountElement        int         `json:"NodeCountElement"`
	NodeCountRecipe         int         `json:"NodeCountRecipe"`
	SearchTimeInMiliseconds int         `json:"SearchTimeInMiliseconds"`
	Steps                   [][3]string `json:"Steps"`
}

func WriteData(pathResult string, elementName string, searchMethod string, isFound bool, steps [][3]string, nodeCountElement int, nodeCountRecipe int, searchTimeMS int) {
	result := SearchResult{
		Element:                 elementName,
		Method:                  searchMethod,
		IsFound:                 isFound,
		NodeCountElement:        nodeCountElement,
		NodeCountRecipe:         nodeCountRecipe,
		SearchTimeInMiliseconds: searchTimeMS,
		Steps:                   steps,
	}

	file, _ := os.Create(pathResult)
	defer file.Close()

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal("Gagal format JSON:", err)
	}

	file.Write(jsonData)

	fmt.Printf("Berhasil menyimpan hasil pencarian [%s] menggunakan metode [%s] di [%s]\n", elementName, searchMethod, pathResult)
}
