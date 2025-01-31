package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Request structure
type GraphRequest struct {
	Edges [][]int `json:"edges"`
	Start int     `json:"start"`
	End   int     `json:"end"`
}

// Response structure
type GraphResponse struct {
	Paths [][]int `json:"paths"`
}

func buildAdjList(edges [][]int) map[int][]int {
	adjList := make(map[int][]int)
	for _, edge := range edges {
		if len(edge) < 2 {
			log.Println("Skipping invalid edge:", edge)
			continue
		}
		u, v := edge[0], edge[1]
		adjList[u] = append(adjList[u], v)
	}
	log.Println("Adjacency List:", adjList)
	return adjList
}

func dfs(adjList map[int][]int, start, end int, path []int, result *[][]int, visited map[int]bool) {
	path = append(path, start) // Add node to path
	visited[start] = true      // Mark node as visited

	if start == end {
		*result = append(*result, append([]int{}, path...))
		log.Println("Found path:", path)
	} else {
		for _, neighbor := range adjList[start] {
			if !visited[neighbor] {
				dfs(adjList, neighbor, end, path, result, visited)
			}
		}
	}

	visited[start] = false
}

func findPathsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request on /find-paths")

	var req GraphRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Println("JSON Decode Error:", err)
		return
	}

	log.Println("Parsed Request:", req)

	adjList := buildAdjList(req.Edges)

	if _, exists := adjList[req.Start]; !exists {
		log.Println("Start node does not exist in adjacency list")
		http.Error(w, "Start node not found", http.StatusNotFound)
		return
	}

	log.Println("Starting DFS Traversal...")

	var result [][]int
	visited := make(map[int]bool) // Track visited nodes
	dfs(adjList, req.Start, req.End, []int{}, &result, visited)

	log.Println("DFS Completed. Paths Found:", result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonData, _ := json.Marshal(GraphResponse{Paths: result})
	w.Write(jsonData)

	log.Println("Response Sent:", string(jsonData))
}

func main() {
	http.HandleFunc("/find-paths", findPathsHandler)
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	fmt.Println("Server running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
