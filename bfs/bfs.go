package bfs

import (
	"fmt"
	"sync"
)

var Reverse map[string][][]string

func IsBaseElement(el string) bool {
	_, found := Reverse[el]
	return !found
}

func BFS(target string) [][]string {
	type State struct {
		Elements []string
		Path     []string
	}

	var results [][]string
	var mu sync.Mutex
	var wg sync.WaitGroup

	queue := make(chan State, 1000000)
	visited := make(map[string]bool)
	visitedMu := sync.Mutex{}

	const workerCount = 8

	queue <- State{Elements: []string{target}, Path: []string{}}

	worker := func() {
		defer wg.Done()
		for state := range queue {
			allBase := true
			for _, el := range state.Elements {
				if !IsBaseElement(el) {
					allBase = false
					break
				}
			}

			if allBase {
				mu.Lock()
				newPath := append([]string{}, state.Path...)
				newPath = append(newPath, state.Elements...)
				results = append(results, newPath)
				mu.Unlock()
				continue
			}

			for i, element := range state.Elements {
				if IsBaseElement(element) {
					continue
				}

				for _, inputs := range Reverse[element] {
					newElements := append([]string{}, state.Elements[:i]...)
					newElements = append(newElements, inputs...)
					newElements = append(newElements, state.Elements[i+1:]...)

					newPath := append([]string{}, state.Path...)
					newPath = append(newPath, element)

					stateKey := fmt.Sprintf("%v|%v", newElements, newPath)

					visitedMu.Lock()
					if !visited[stateKey] {
						visited[stateKey] = true
						visitedMu.Unlock()

						queue <- State{
							Elements: newElements,
							Path:     newPath,
						}
					} else {
						visitedMu.Unlock()
					}
				}
			}
		}
	}

	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker()
	}

	go func() {
		wg.Wait()
		close(queue)
	}()

	wg.Wait()

	return results
}

func FindShortestPath(target string) []string {
	allPaths := BFS(target)
	if len(allPaths) == 0 {
		return nil
	}

	shortest := allPaths[0]
	for _, path := range allPaths {
		if len(path) < len(shortest) {
			shortest = path
		}
	}
	return shortest
}

// // Testing: bfs_test.go

// package bfs_test

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"testing"

// 	"littlealchemy/bfs"
// )

// type TreeNode struct {
// 	Name     string
// 	Children []*TreeNode
// }

// func loadRecipes() error {
// 	path := filepath.Join("..", "data", "recipe.json")
// 	file, err := os.ReadFile(path)
// 	if err != nil {
// 		return err
// 	}

// 	var recipes []struct {
// 		Output string     `json:"Output"`
// 		Inputs [][]string `json:"Inputs"`
// 	}
// 	if err := json.Unmarshal(file, &recipes); err != nil {
// 		return err
// 	}

// 	bfs.Reverse = make(map[string][][]string)
// 	for _, recipe := range recipes {
// 		bfs.Reverse[recipe.Output] = recipe.Inputs
// 	}
// 	return nil
// }

// func buildSpecificTree(element string, recipeIndex int) *TreeNode {
// 	node := &TreeNode{Name: element}

// 	if inputsList, ok := bfs.Reverse[element]; ok {
// 		if recipeIndex < len(inputsList) {
// 			inputs := inputsList[recipeIndex]
// 			for _, input := range inputs {
// 				child := buildSpecificTree(input, 0)
// 				node.Children = append(node.Children, child)
// 			}
// 		}
// 	}

// 	return node
// }

// func printTree(node *TreeNode, prefix string, isTail bool) {
// 	if node == nil {
// 		return
// 	}
// 	branch := "├── "
// 	if isTail {
// 		branch = "└── "
// 	}
// 	fmt.Println(prefix + branch + node.Name)

// 	for i := 0; i < len(node.Children); i++ {
// 		nextPrefix := prefix
// 		if isTail {
// 			nextPrefix += "    "
// 		} else {
// 			nextPrefix += "│   "
// 		}
// 		printTree(node.Children[i], nextPrefix, i == len(node.Children)-1)
// 	}
// }

// func printPath(path []string) {
// 	for i, element := range path {
// 		if i > 0 {
// 			fmt.Print(" -> ")
// 		}
// 		fmt.Print(element)
// 	}
// 	fmt.Println()
// }

// func TestBFSTree(t *testing.T) {
// 	err := loadRecipes()
// 	if err != nil {
// 		t.Fatalf("Failed to load recipes: %v", err)
// 	}

// 	target := "result"

// 	if recipes, ok := bfs.Reverse[target]; ok {
// 		for i, recipe := range recipes {
// 			fmt.Printf("\nHasil %d (Resep: %v):\n", i+1, recipe)
// 			tree := &TreeNode{Name: target}
// 			for _, input := range recipe {
// 				child := buildSpecificTree(input, 0)
// 				tree.Children = append(tree.Children, child)
// 			}
// 			printTree(tree, "", true)
// 		}
// 	} else {
// 		t.Fatalf("No recipes found for %s", target)
// 	}
// }

// func TestBFSAllPaths(t *testing.T) {
// 	err := loadRecipes()
// 	if err != nil {
// 		t.Fatalf("Failed to load recipes: %v", err)
// 	}

// 	target := "result"
// 	paths := bfs.BFS(target)

// 	fmt.Printf("\nAll paths to create %s:\n", target)
// 	for i, path := range paths {
// 		fmt.Printf("Path %d: ", i+1)
// 		printPath(path)
// 	}
// }

// func TestBFSShortestPath(t *testing.T) {
// 	err := loadRecipes()
// 	if err != nil {
// 		t.Fatalf("Failed to load recipes: %v", err)
// 	}

// 	target := "result"
// 	shortestPath := bfs.FindShortestPath(target)

// 	fmt.Printf("\nShortest path to create %s:\n", target)
// 	printPath(shortestPath)
// }
