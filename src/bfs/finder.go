package bfs

import (
	"fmt"
	"log"
	"time"
)

func BFS(elementName string, pathResult string, isMultiple bool, nRecipe int) {
	if isMultiple {
		helperBFSMultiple(elementName, pathResult, nRecipe)
	} else {
		helperBFSSingle(elementName, pathResult)
	}
}

/* ==================================================================================================================================== */

func helperBFSSingle(elementName string, pathResult string) {
	recipeMap, _, err := LoadData("./data/recipe.json")
	if err != nil {
		log.Fatal("Gagal memuat file JSON:", err)
	}

	start := time.Now()
	ok, steps := BFSSingle(elementName, recipeMap)
	elapsed := time.Since(start)

	// Untuk kompatibilitas dengan fungsi WriteData yang ada
	nodeCountElement := len(steps) * 2 // Perkiraan sederhana
	nodeCountRecipe := len(steps)

	WriteData(pathResult, elementName, "BFS", ok, steps, nodeCountElement, nodeCountRecipe, int(elapsed.Milliseconds()))

	if !ok {
		fmt.Println("Tidak dapat membuat", elementName)
		return
	}

	fmt.Println("Total step: ", len(steps))
	fmt.Printf("Searching Time: %d ms\n", elapsed.Milliseconds())
}

// Implementasi BFS murni tanpa optimasi tier
func BFSSingle(elementName string, recipeMap map[string][][2]string) (bool, [][3]string) {
	baseElements := map[string]bool{"Air": true, "Water": true, "Fire": true, "Earth": true, "Time": true}

	// Queue berisi elemen yang perlu diproses dan path yang sudah dilalui
	queue := []struct {
		element string
		path    [][3]string
	}{
		{elementName, [][3]string{}},
	}

	visited := make(map[string]bool)
	visited[elementName] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if baseElements[current.element] {
			continue
		}

		recipes, exists := recipeMap[current.element]
		if !exists {
			continue
		}

		for _, recipe := range recipes {
			nameA, nameB := recipe[0], recipe[1]
			newPath := append(current.path, [3]string{nameA, nameB, current.element})

			// Jika kedua komponen adalah base element
			if baseElements[nameA] && baseElements[nameB] {
				return true, newPath
			}

			// Proses komponen A jika belum dikunjungi
			if !visited[nameA] {
				visited[nameA] = true
				queue = append(queue, struct {
					element string
					path    [][3]string
				}{nameA, newPath})
			}

			// Proses komponen B jika belum dikunjungi
			if !visited[nameB] {
				visited[nameB] = true
				queue = append(queue, struct {
					element string
					path    [][3]string
				}{nameB, newPath})
			}
		}
	}

	return false, nil
}

/* ===================================================================================================================== */

func helperBFSMultiple(elementName string, pathResult string, nPath int) {
	recipeMap, _, err := LoadData("./data/recipe.json")
	if err != nil {
		log.Fatal("Gagal memuat file JSON:", err)
	}

	start := time.Now()
	results := BFSMultiple(elementName, recipeMap, nPath)
	elapsed := time.Since(start)

	// Untuk kompatibilitas dengan fungsi WriteData yang ada
	nodeCountElement := 0
	nodeCountRecipe := 0
	for _, path := range results {
		nodeCountElement += len(path) * 2
		nodeCountRecipe += len(path)
	}

	isSatisfied := len(results) == nPath
	WriteDataMultiplePaths(pathResult, elementName, "BFS", isSatisfied, results, nodeCountElement, nodeCountRecipe, int(elapsed.Milliseconds()), nPath)

	fmt.Println("Target element: ", elementName)
	fmt.Print("Total founded recipe: ", len(results), "\n")
	fmt.Printf("Searching Time: %d ms\n", elapsed.Milliseconds())
}

// Implementasi BFS multiple yang lebih murni
func BFSMultiple(
	elementName string,
	recipeMap map[string][][2]string,
	maxPaths int,
) []RecipePath {
	baseElements := map[string]bool{"Air": true, "Water": true, "Fire": true, "Earth": true, "Time": true}

	type queueItem struct {
		element string
		path    RecipePath
	}

	queue := []queueItem{{elementName, RecipePath{}}}
	results := []RecipePath{}
	visited := make(map[string]bool)
	visited[elementName] = true

	for len(queue) > 0 && len(results) < maxPaths {
		current := queue[0]
		queue = queue[1:]

		if baseElements[current.element] {
			continue
		}

		recipes, exists := recipeMap[current.element]
		if !exists {
			continue
		}

		for _, recipe := range recipes {
			nameA, nameB := recipe[0], recipe[1]
			newStep := [3]string{nameA, nameB, current.element}
			newPath := append(current.path, newStep)

			if baseElements[nameA] && baseElements[nameB] {
				results = append(results, newPath)
				if len(results) >= maxPaths {
					return results
				}
				continue
			}

			if !visited[nameA] {
				visited[nameA] = true
				queue = append(queue, queueItem{nameA, newPath})
			}

			if !visited[nameB] {
				visited[nameB] = true
				queue = append(queue, queueItem{nameB, newPath})
			}
		}
	}

	return results
}
