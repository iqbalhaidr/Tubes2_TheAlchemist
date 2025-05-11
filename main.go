package main

import (
	"littlealchemy/dfs"
)

func main() {
	// dfs.DFS("Bank", "./data/result.json")
	dfs.DFSN("Gold", "./data/result.json", 5)
}
