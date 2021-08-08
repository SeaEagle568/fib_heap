package dijkstra

import (
	"github.com/SeaEagle568/fib_heap/pkg/fib_heap"
)

func DijkstraHeap(g Graph, start int) (dist []int64, p []Node) {
	dist = make([]int64, g.Size+1)
	p = make([]int, g.Size+1)
	for i, _ := range dist {
		dist[i] = Inf
	}
	dist[start] = 0
	heap := fib_heap.FibonacciHeap{}
	heap.Insert(IntPair{start, 0})
	for !heap.IsEmpty() {
		var cur IntPair
		tmp, _ := heap.ExtractMin()
		cur = tmp.(IntPair)
		v := cur.first

		for _, pair := range g.edges[v] {
			to := pair.first
			cost := pair.second
			if dist[to] > dist[v]+cost {
				heap.DeleteKey(IntPair{to, dist[to]})
				dist[to] = dist[v] + cost
				p[to] = v
				heap.Insert(IntPair{to, dist[to]})
			}

		}
	}

	return
}

func ShortestPathHeap(g Graph, start int, finish int) (dist int64, path []Node) {
	d, pr := DijkstraHeap(g, start)
	dist = d[finish]
	if dist == Inf {
		return
	}
	path = make([]Node, 0)
	path = append(path, finish)
	for start != finish {
		finish = pr[finish]
		path = append(path, finish)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return
}
