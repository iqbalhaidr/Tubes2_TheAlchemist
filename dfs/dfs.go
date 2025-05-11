package dfs

import (
	"fmt"
	"log"
	"time"
)

func DFS(elementName string, pathResult string) {
	recipeMap, tierMap, err := LoadData("./data/recipe.json")
	if err != nil {
		log.Fatal("Gagal memuat file JSON:", err)
	}

	target := elementName
	var nodeCountElement int = 0
	var nodeCountRecipe int = 0
	visitedMap := make(map[string]bool)

	start := time.Now()
	ok, steps := dfsSingle(target, recipeMap, tierMap, &nodeCountElement, &nodeCountRecipe, visitedMap)
	elapsed := time.Since(start)

	WriteData(pathResult, target, "DFS", ok, steps, nodeCountElement, nodeCountRecipe, int(elapsed.Milliseconds()))

	if !ok {
		fmt.Println("Tidak dapat membuat", target)
		return
	}

	fmt.Println("nodeCountElement: ", nodeCountElement)
	fmt.Println("nodeCountRecipe: ", nodeCountRecipe)
	fmt.Println("Total step: ", len(steps))
	fmt.Printf("Searching Time: %d ms\n", elapsed.Milliseconds())
}

// DFS mencari langkah-langkah untuk membentuk elemen
func dfsSingle(elementName string, recipeMap map[string][][2]string, tierMap map[string]int, nodeCountElement *int, nodeCountRecipe *int, visitedMap map[string]bool) (bool, [][3]string) {
	(*nodeCountElement)++
	base := map[string]bool{"Air": true, "Water": true, "Fire": true, "Earth": true, "Time": true}
	if base[elementName] {
		return true, [][3]string{}
	}

	if visitedMap[elementName] {
		return true, [][3]string{}
	}

	elementTier := tierMap[elementName]

	for _, recipe := range recipeMap[elementName] {
		nameA := recipe[0]
		nameB := recipe[1]
		tierA := tierMap[nameA]
		tierB := tierMap[nameB]

		if tierA >= elementTier || tierB >= elementTier {
			continue
		}

		okA, stepsA := dfsSingle(nameA, recipeMap, tierMap, nodeCountElement, nodeCountRecipe, visitedMap)
		if !okA {
			continue
		}

		okB, stepsB := dfsSingle(nameB, recipeMap, tierMap, nodeCountElement, nodeCountRecipe, visitedMap)
		if !okB {
			continue
		}

		combined := append(append(stepsA, stepsB...), [3]string{nameA, nameB, elementName})
		(*nodeCountRecipe)++
		visitedMap[elementName] = true
		return true, combined
	}

	return false, nil
}
