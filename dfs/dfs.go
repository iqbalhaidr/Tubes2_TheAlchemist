package dfs

func dfs(target string, path []string, visited map[string]bool) [][]string {
    if isBaseElement(target) {
        return [][]string{{target}}
    }

    var results [][]string
    for _, inputs := range reverse[target] {
        if visited[target] {
            continue
        }
        visited[target] = true

        leftPaths := dfs(inputs[0], path, visited)
        rightPaths := dfs(inputs[1], path, visited)

        for _, l := range leftPaths {
            for _, r := range rightPaths {
                results = append(results, append(append([]string{}, l...), r..., target))
            }
        }
        visited[target] = false
    }
    return results
}
