package bfs

import "fmt"

var Reverse map[string][][]string

func IsBaseElement(el string) bool {
	_, found := Reverse[el]
	return !found
}

func Bfs(target string) [][]string {
	type State struct {
		Elements []string
		Path     []string
	}

	var results [][]string
	queue := []State{{Elements: []string{target}, Path: []string{}}}
	visited := make(map[string]bool)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		allBase := true
		for _, el := range curr.Elements {
			if !IsBaseElement(el) {
				allBase = false
				break
			}
		}
		if allBase {
			newPath := append([]string{}, curr.Path...)
			newPath = append(newPath, curr.Elements...)
			results = append(results, newPath)
			continue
		}

		for i, element := range curr.Elements {
			if IsBaseElement(element) {
				continue
			}
			for _, inputs := range Reverse[element] {
				newElements := append([]string{}, curr.Elements[:i]...)
				newElements = append(newElements, inputs...)
				newElements = append(newElements, curr.Elements[i+1:]...)

				newPath := append([]string{}, curr.Path...)
				newPath = append(newPath, element)
				newElementsStr := fmt.Sprint(newElements)
				if !visited[newElementsStr] {
					queue = append(queue, State{Elements: newElements, Path: newPath})
					visited[newElementsStr] = true
				}
			}
			break
		}
	}

	return results
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

// func TestBFSTree(t *testing.T) {
// 	err := loadRecipes()
// 	if err != nil {
// 		t.Fatalf("Failed to load recipes: %v", err)
// 	}

// 	target := "final"

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
