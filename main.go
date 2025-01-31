package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GraphRequest struct {
	Edges [][]int `json:"edges"`
	Start int     `json:"start"`
	End   int     `json:"end"`
}

type GraphResponse struct {
	Paths [][]int `json:"paths"`
}

func buildAdjList(edges [][]int) map[int][]int {
	adjList := make(map[int][]int)

	for _, edge := range edges {
		u, v := edge[0], edge[1]
		adjList[u] = append(adjList[u], v)

	}

	return adjList

}

//DFS

func dfs(adjList map[int][]int, start, end int, path []int, result *[][]int, visited map[int]bool) {
	path = append(path, start)
	log.Println("DFS Visiting node : ", start, "path: ", path)
	if start == end {
		*result = append(*result, append([]int{}, path...))
		return
	}

	//travel through the neighbours

	for _, neighbor := range adjList[start] {
		if !visited[neighbor] {
			dfs(adjList, neighbor, end, path, result, visited)
		}
	}

}

//handler for Endpoint: /find-paths

func findPathsHandler(w http.ResponseWriter, r *http.Request) {
	var req GraphRequest
	//To handle any error if the request is invalid
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}
	return

	adjList := buildAdjList(req.Edges)

	//edge cases
	if _, exists := adjList[req.Start]; !exists {
		fmt.Print(exists)
		json.NewEncoder(w).Encode(GraphResponse{Paths: [][]int{}})

		return
	}

	var result [][]int
	visited := make(map[int]bool)

	//do the dfs
	dfs(adjList, req.Start, req.End, []int{}, &result, visited)
	fmt.Print(adjList)
	//for response with path
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GraphResponse{Paths: result})

}

func main() {
	http.HandleFunc("/find-paths", findPathsHandler)
	//logs
	fmt.Println("Server running on port 8081 ... ")
	log.Fatal(http.ListenAndServe("192.168.122.238:8081", nil))
}
