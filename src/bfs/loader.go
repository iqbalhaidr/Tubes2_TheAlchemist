package bfs

import (
	"encoding/json"
	"os"
)

// RecipeEntry menyimpan satu resep dengan output dan bahan-bahannya
type RecipeEntry struct {
	Output string      `json:"Output"`
	Inputs [][2]string `json:"Inputs"`
	Tier   int         `json:"Tier"`
}

// LoadData membaca file JSON dan mengembalikan peta resep dan peta tier
func LoadData(filename string) (map[string][][2]string, map[string]int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	var entries []RecipeEntry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, nil, err
	}

	recipeMap := make(map[string][][2]string)
	tierMap := make(map[string]int)

	for _, entry := range entries {
		recipeMap[entry.Output] = entry.Inputs
		tierMap[entry.Output] = entry.Tier
	}

	return recipeMap, tierMap, nil
}
