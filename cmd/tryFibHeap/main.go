package main

import (
	"bufio"
	"fmt"
	"github.com/SeaEagle568/fib_heap/pkg/dijkstra"
	"github.com/SeaEagle568/fib_heap/pkg/fib_heap"
	"io"
	"os"
	"strconv"
)

/**
Testing two Dijkstra implementations on massive testcase
*/

func main() {

	file, _ := os.Open("test/case1.in")
	input, _ := ReadInts(file)
	var n, m int
	n = input[0]
	m = input[1]
	fmt.Println(n, m)
	count := 2
	g := dijkstra.Graph{Size: n}
	for i := 0; i < m; i++ {
		var from, to, cost int
		from = input[count]
		to = input[count+1]
		cost = input[count+2]
		count += 3
		g.AddEdge(dijkstra.Edge{From: from, To: to, Cost: int64(cost)})
	}
	dist, path := dijkstra.ShortestPathSet(g, 1, n)
	fmt.Printf("Set dijkstra answer: %v\n", dist)
	if dist == -dijkstra.Inf {
		fmt.Println(-1)
	} else {
		for _, val := range path {
			fmt.Printf("%v ", val)
		}
	}
	fmt.Println()

	dist, path = dijkstra.ShortestPathHeap(g, 1, n)
	fmt.Printf("Heap dijkstra answer: %v\n", dist)
	if dist == -dijkstra.Inf {
		fmt.Println(-1)
	} else {
		for _, val := range path {
			fmt.Printf("%v ", val)
		}
	}
	fmt.Println("\n\n", fib_heap.Timer, fib_heap.Timer2)
}

func ReadInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var result []int
	for scanner.Scan() {
		x, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return result, err
		}
		result = append(result, x)
	}
	return result, scanner.Err()
}
