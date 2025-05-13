package bfs

import (
	"fmt"
	"log"
	"sync"
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
	recipeMap, tierMap, err := LoadData("./data/recipe.json")
	if err != nil {
		log.Fatal("Gagal memuat file JSON:", err)
	}

	target := elementName
	var nodeCountElement int = 0
	var nodeCountRecipe int = 0
	memo := make(map[string][]RecipePath)

	start := time.Now()
	// Menggunakan bfsMultiple dengan nPath = 1
	hasil := bfsMultiple(target, recipeMap, tierMap, memo, &nodeCountElement, &nodeCountRecipe, 1)
	elapsed := time.Since(start)

	var steps [][3]string
	ok := len(hasil) > 0
	if ok {
		// Mengonversi hasil ke format penulisan
		steps = make([][3]string, len(hasil[0]))
		for i, step := range hasil[0] {
			steps[i] = [3]string(step)
		}
	}

	WriteData(pathResult, target, "BFS", ok, steps, nodeCountElement, nodeCountRecipe, int(elapsed.Milliseconds()))

	if !ok {
		fmt.Println("Tidak dapat membuat", target)
		return
	}

	fmt.Println("nodeCountElement: ", nodeCountElement)
	fmt.Println("nodeCountRecipe: ", nodeCountRecipe)
	fmt.Println("Total step: ", len(steps))
	fmt.Printf("Searching Time: %d ms\n", elapsed.Milliseconds())
}

/* ===================================================================================================================== */

type Step [3]string
type RecipePath []Step

var baseElements = map[string]bool{
	"Air": true, "Water": true, "Fire": true, "Earth": true, "Time": true,
}

var memoLock sync.RWMutex

func helperBFSMultiple(elementName string, pathResult string, nPath int) {
	recipeMap, tierMap, err := LoadData("./data/recipe.json")
	if err != nil {
		log.Fatal("Gagal memuat file JSON:", err)
	}

	target := elementName
	var nodeCountElement int = 0
	var nodeCountRecipe int = 0
	memo := make(map[string][]RecipePath)

	start := time.Now()
	hasil := bfsMultiple(target, recipeMap, tierMap, memo, &nodeCountElement, &nodeCountRecipe, nPath)
	elapsed := time.Since(start)

	var isSatisfied bool = len(hasil) == nPath
	WriteDataMultiplePaths(pathResult, target, "BFS", isSatisfied, hasil, nodeCountElement, nodeCountRecipe, int(elapsed.Milliseconds()), nPath)

	fmt.Println("Target element: ", elementName)
	fmt.Print("Total founded recipe: ", len(hasil), "\n")
	fmt.Println("nodeCountElement: ", nodeCountElement)
	fmt.Println("nodeCountRecipe: ", nodeCountRecipe)
	fmt.Printf("Searching Time: %d ms\n", elapsed.Milliseconds())
}

func bfsMultiple(
	elementName string,
	recipeMap map[string][][2]string,
	tierMap map[string]int,
	memo map[string][]RecipePath,
	nodeCountElement *int,
	nodeCountRecipe *int,
	maxPaths int,
) []RecipePath {
	(*nodeCountElement)++

	// Mengembalikan jalur kosong jika elemen dasar
	if baseElements[elementName] {
		return []RecipePath{{}}
	}

	// Memeriksa memo agar tidak menghitung ulang
	memoLock.RLock()
	if paths, found := memo[elementName]; found {
		memoLock.RUnlock()
		return paths
	}
	memoLock.RUnlock()

	type QueueItem struct {
		element string
		paths   []RecipePath
	}

	queue := []QueueItem{{element: elementName, paths: []RecipePath{}}}
	solutions := make(map[string][]RecipePath)
	elementTier := tierMap[elementName]

	for len(queue) > 0 && len(solutions[elementName]) < maxPaths {
		current := queue[0]
		queue = queue[1:]
		(*nodeCountElement)++

		if baseElements[current.element] {
			solutions[current.element] = []RecipePath{{}}
			continue
		}

		// Mencoba semua resep untuk elemen
		for _, recipe := range recipeMap[current.element] {
			nameA, nameB := recipe[0], recipe[1]
			tierA, tierB := tierMap[nameA], tierMap[nameB]

			// Menghindari siklus tak berujung dengan tier lebih tinggi atau sama
			if tierA >= elementTier || tierB >= elementTier {
				continue
			}

			(*nodeCountRecipe)++

			pathsA, hasA := solutions[nameA]
			if !hasA {
				pathsA = bfsMultiple(nameA, recipeMap, tierMap, memo, nodeCountElement, nodeCountRecipe, maxPaths)
				solutions[nameA] = pathsA
			}

			pathsB, hasB := solutions[nameB]
			if !hasB {
				pathsB = bfsMultiple(nameB, recipeMap, tierMap, memo, nodeCountElement, nodeCountRecipe, maxPaths)
				solutions[nameB] = pathsB
			}

			var newPaths []RecipePath
			for _, pathA := range pathsA {
				for _, pathB := range pathsB {
					// Menggabungkan path dari A dan B serta langkah membuat elemen sekarang
					combined := make(RecipePath, 0, len(pathA)+len(pathB)+1)
					combined = append(combined, pathA...)
					combined = append(combined, pathB...)
					combined = append(combined, Step{nameA, nameB, current.element})
					newPaths = append(newPaths, combined)

					// Menghentikan jika sudah cukup banyak jalur
					if current.element == elementName && len(newPaths) >= maxPaths {
						break
					}
				}
				if current.element == elementName && len(newPaths) >= maxPaths {
					break
				}
			}

			if len(newPaths) > 0 {
				if _, exists := solutions[current.element]; !exists {
					solutions[current.element] = newPaths
				} else {
					solutions[current.element] = append(solutions[current.element], newPaths...)
				}

				if current.element == elementName && len(solutions[current.element]) >= maxPaths {
					solutions[current.element] = solutions[current.element][:maxPaths]
					break
				}
			}
		}
	}

	// Menyimpan ke memo untuk menghindari perhitungan ulang
	memoLock.Lock()
	memo[elementName] = solutions[elementName]
	memoLock.Unlock()

	return solutions[elementName]
}
