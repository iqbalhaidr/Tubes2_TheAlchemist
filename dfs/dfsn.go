package dfs

import (
	"fmt"
	"log"
	"time"
)

type Step [3]string    // Representasi satu langkah: A + B => C
type RecipePath []Step // Satu jalur lengkap pembentukan elemen

func DFSN(elementName string, pathResult string, nPath int) {
	recipeMap, tierMap, err := LoadData("./data/recipe.json")
	if err != nil {
		log.Fatal("Gagal memuat file JSON:", err)
	}

	target := elementName
	var nodeCountElement int = 0
	var nodeCountRecipe int = 0
	memo := make(map[string][]RecipePath)

	start := time.Now()
	hasil := DfsN(target, recipeMap, tierMap, memo, &nodeCountElement, nPath)
	elapsed := time.Since(start)

	WriteDataMultiplePaths(pathResult, target, "DFS", true, hasil, nodeCountElement, int(elapsed.Milliseconds()))

	// if !ok {
	// 	fmt.Println("Tidak dapat membuat", target)
	// 	return
	// }

	fmt.Println("nodeCountElement: ", nodeCountElement)
	fmt.Println("nodeCountRecipe: ", nodeCountRecipe)
	// fmt.Println("Total step: ", len(steps))
	fmt.Printf("Searching Time: %d ms\n", elapsed.Milliseconds())
}

func DfsN(
	elementName string,
	recipeMap map[string][][2]string,
	tierMap map[string]int,
	memo map[string][]RecipePath,
	nodeCountElement *int,
	maxPaths int,
) []RecipePath {
	(*nodeCountElement)++

	base := map[string]bool{"Air": true, "Water": true, "Fire": true, "Earth": true, "Time": true}
	if base[elementName] {
		return []RecipePath{{}} // Base element: satu path kosong
	}

	if paths, found := memo[elementName]; found {
		return paths
	}

	elementTier := tierMap[elementName]
	var allPaths []RecipePath

	for _, recipe := range recipeMap[elementName] {
		nameA, nameB := recipe[0], recipe[1]
		tierA, tierB := tierMap[nameA], tierMap[nameB]

		if tierA >= elementTier || tierB >= elementTier {
			continue
		}

		pathsA := DfsN(nameA, recipeMap, tierMap, memo, nodeCountElement, maxPaths)
		pathsB := DfsN(nameB, recipeMap, tierMap, memo, nodeCountElement, maxPaths)

		for _, pathA := range pathsA {
			for _, pathB := range pathsB {
				combined := append([]Step{}, pathA...)
				combined = append(combined, pathB...)
				combined = append(combined, Step{nameA, nameB, elementName})
				allPaths = append(allPaths, combined)

				if len(allPaths) >= maxPaths {
					memo[elementName] = allPaths
					return allPaths
				}
			}
		}
	}

	memo[elementName] = allPaths
	return allPaths
}
