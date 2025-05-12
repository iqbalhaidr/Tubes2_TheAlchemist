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

	start := time.Now()
	ok, steps := bfsSingle(target, recipeMap, tierMap, &nodeCountElement, &nodeCountRecipe)
	elapsed := time.Since(start)

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

func bfsSingle(elementName string, recipeMap map[string][][2]string, tierMap map[string]int, nodeCountElement *int, nodeCountRecipe *int) (bool, [][3]string) {
	base := map[string]bool{"Air": true, "Water": true, "Fire": true, "Earth": true, "Time": true}
	if base[elementName] {
		return true, [][3]string{}
	}

	type QueueItem struct {
		element string
		path    [][3]string
	}

	queue := []QueueItem{{element: elementName, path: [][3]string{}}}
	visited := make(map[string]bool)
	parent := make(map[string][2]string)
	solutionPath := make(map[string][][3]string)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		(*nodeCountElement)++

		if visited[current.element] {
			continue
		}
		visited[current.element] = true

		if base[current.element] {
			solutionPath[current.element] = [][3]string{}
			continue
		}

		elementTier := tierMap[current.element]
		foundSolution := false

		for _, recipe := range recipeMap[current.element] {
			nameA := recipe[0]
			nameB := recipe[1]
			tierA := tierMap[nameA]
			tierB := tierMap[nameB]

			if tierA >= elementTier || tierB >= elementTier {
				continue
			}

			(*nodeCountRecipe)++

			if _, existsA := solutionPath[nameA]; !existsA {
				queue = append(queue, QueueItem{element: nameA, path: append(current.path, [3]string{nameA, nameB, current.element})})
			}
			if _, existsB := solutionPath[nameB]; !existsB {
				queue = append(queue, QueueItem{element: nameB, path: append(current.path, [3]string{nameA, nameB, current.element})})
			}

			if solutionA, okA := solutionPath[nameA]; okA {
				if solutionB, okB := solutionPath[nameB]; okB {
					combined := append(append(solutionA, solutionB...), [3]string{nameA, nameB, current.element})
					solutionPath[current.element] = combined
					foundSolution = true
					if current.element == elementName {
						return true, combined
					}
					break
				}
			}
		}

		if !foundSolution {
			parent[current.element] = [2]string{}
		}
	}

	if path, ok := solutionPath[elementName]; ok {
		return true, path
	}
	return false, nil
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

	if baseElements[elementName] {
		return []RecipePath{{}}
	}

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

		for _, recipe := range recipeMap[current.element] {
			nameA, nameB := recipe[0], recipe[1]
			tierA, tierB := tierMap[nameA], tierMap[nameB]

			if tierA >= elementTier || tierB >= elementTier {
				continue
			}

			(*nodeCountRecipe)++

			// Check if we have solutions for the components
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
					combined := make(RecipePath, 0, len(pathA)+len(pathB)+1)
					combined = append(combined, pathA...)
					combined = append(combined, pathB...)
					combined = append(combined, Step{nameA, nameB, current.element})
					newPaths = append(newPaths, combined)

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

	memoLock.Lock()
	memo[elementName] = solutions[elementName]
	memoLock.Unlock()

	return solutions[elementName]
}
