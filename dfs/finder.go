package dfs

import (
	"fmt"
	"log"
	"sync"
	"time"
)

/*
YANG DIPANGGIL YANG INI
elementName = nama element yang ingin dicari resepnya
pathResult = tempat menyimpan file hasil pencarian, berupa relative path dari parent directory i.e. "./data/result.json"
isMultiple = true jika ingin mencari multiple recipe dari sebuah elemen, vice versa
nRecipe = banyaknya resep berbeda yang ingin dicari, jika isMultiple false maka isi 1
*/
func DFS(elementName string, pathResult string, isMultiple bool, nRecipe int) {
	if isMultiple {
		helperDFSMultiple(elementName, pathResult, nRecipe)
	} else {
		helperDFSSingle(elementName, pathResult)
	}
}

/* ==================================================================================================================================== */

func helperDFSSingle(elementName string, pathResult string) {
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

/* ===================================================================================================================== */

type Step [3]string    // Representasi satu langkah: A + B => C
type RecipePath []Step // Satu jalur lengkap pembentukan elemen
var baseElements = map[string]bool{
	"Air": true, "Water": true, "Fire": true, "Earth": true, "Time": true,
}

// mutex global untuk memo agar thread-safe
var memoLock sync.RWMutex

func helperDFSMultiple(elementName string, pathResult string, nPath int) {
	recipeMap, tierMap, err := LoadData("./data/recipe.json")
	if err != nil {
		log.Fatal("Gagal memuat file JSON:", err)
	}

	target := elementName
	var nodeCountElement int = 0
	var nodeCountRecipe int = 0
	memo := make(map[string][]RecipePath)

	start := time.Now()
	hasil := dfsMultiple(target, recipeMap, tierMap, memo, &nodeCountElement, &nodeCountRecipe, nPath)
	elapsed := time.Since(start)

	var isSatisfied bool = len(hasil) == nPath
	WriteDataMultiplePaths(pathResult, target, "DFS", isSatisfied, hasil, nodeCountElement, nodeCountRecipe, int(elapsed.Milliseconds()), nPath)

	fmt.Println("Target element: ", elementName)
	fmt.Print("Total founded recipe: ", len(hasil), "\n")
	fmt.Println("nodeCountElement: ", nodeCountElement)
	fmt.Println("nodeCountRecipe: ", nodeCountRecipe)
	fmt.Printf("Searching Time: %d ms\n", elapsed.Milliseconds())
}

func dfsMultiple(
	elementName string,
	recipeMap map[string][][2]string,
	tierMap map[string]int,
	memo map[string][]RecipePath,
	nodeCountElement *int,
	nodeCountRecipe *int,
	maxPaths int,
) []RecipePath {
	(*nodeCountElement)++

	if baseElements[elementName] {
		return []RecipePath{{}} // Base element: satu path kosong
	}

	// Cek di memo dengan read lock
	memoLock.RLock()
	if paths, found := memo[elementName]; found {
		memoLock.RUnlock()
		return paths
	}
	memoLock.RUnlock()

	elementTier := tierMap[elementName]
	var allPaths []RecipePath

	for _, recipe := range recipeMap[elementName] {
		nameA, nameB := recipe[0], recipe[1]
		tierA, tierB := tierMap[nameA], tierMap[nameB]

		if tierA >= elementTier || tierB >= elementTier {
			continue
		}

		// Paralelkan pencarian nameA dan nameB
		chA := make(chan []RecipePath, 1)
		chB := make(chan []RecipePath, 1)

		go func() {
			chA <- getPathsConcurrent(nameA, recipeMap, tierMap, memo, nodeCountElement, nodeCountRecipe, maxPaths)
		}()

		go func() {
			chB <- getPathsConcurrent(nameB, recipeMap, tierMap, memo, nodeCountElement, nodeCountRecipe, maxPaths)
		}()

		pathsA := <-chA
		pathsB := <-chB

		(*nodeCountRecipe)++

		for _, pathA := range pathsA {
			for _, pathB := range pathsB {
				combined := append([]Step{}, pathA...)
				combined = append(combined, pathB...)
				combined = append(combined, Step{nameA, nameB, elementName})
				allPaths = append(allPaths, combined)

				if len(allPaths) >= maxPaths {
					// Simpan ke memo dengan write lock
					memoLock.Lock()
					memo[elementName] = allPaths
					memoLock.Unlock()
					return allPaths
				}
			}
		}
	}

	// Simpan hasil akhir ke memo
	memoLock.Lock()
	memo[elementName] = allPaths
	memoLock.Unlock()
	return allPaths
}

func getPathsConcurrent(
	name string,
	recipeMap map[string][][2]string,
	tierMap map[string]int,
	memo map[string][]RecipePath,
	nodeCountElement *int,
	nodeCountRecipe *int,
	maxPaths int,
) []RecipePath {
	memoLock.RLock()
	if paths, found := memo[name]; found {
		memoLock.RUnlock()
		return paths
	}
	memoLock.RUnlock()

	result := dfsMultiple(name, recipeMap, tierMap, memo, nodeCountElement, nodeCountRecipe, maxPaths)

	memoLock.Lock()
	memo[name] = result
	memoLock.Unlock()

	return result
}
