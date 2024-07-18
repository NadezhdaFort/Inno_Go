package main

import (
	"fmt"
	"math"
)

const Inf = math.MaxInt

type GraphMtx struct {
	adjMtx [][]int
}

func NewGraphMtx(mtx [][]int) *GraphMtx {
	return &GraphMtx{adjMtx: mtx}
}

type emptyVal struct{}

func (g *GraphMtx) BFS(startVertex int) []int {
	path := make([]int, 0, len(g.adjMtx))
	visited := make(map[int]emptyVal)
	var queue []int

	queue = append(queue, startVertex)
	visited[startVertex] = emptyVal{}

	for len(queue) != 0 {
		currentVertex := queue[0]
		queue = queue[1:]
		path = append(path, currentVertex)
		for idx, cost := range g.adjMtx[currentVertex] {
			if _, ok := visited[idx]; !ok && cost != Inf {
				queue = append(queue, idx)
				visited[idx] = emptyVal{}
			}
		}
	}
	return path
}

func main() {
	// двунаправленный граф
	matrix := [][]int{
		{0, 2, Inf, 6, Inf},
		{2, 0, 3, 8, 5},
		{Inf, 3, 0, Inf, 7},
		{6, 8, Inf, 0, 9},
		{Inf, 5, 7, 9, 0},
	}
	// однонаправленный граф
	matrix2 := [][]int{
		{0, 3, 1, Inf},
		{Inf, 0, 7, Inf},
		{Inf, Inf, 0, 2},
		{Inf, Inf, Inf, 0},
	}
	graph := NewGraphMtx(matrix)
	graph2 := NewGraphMtx(matrix2)
	fmt.Println(graph.BFS(4))  // [4 1 2 3 0]
	fmt.Println(graph.BFS(2))  // [2 1 4 0 3]
	fmt.Println(graph2.BFS(0)) // [0 1 2 3]
	fmt.Println(graph2.BFS(1)) // [1 2 3]
}
