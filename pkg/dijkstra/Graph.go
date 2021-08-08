package dijkstra

import "github.com/SeaEagle568/fib_heap/pkg/fib_heap"

type Node = int

type IntPair struct {
	first  int
	second int64
}

func (i IntPair) LessThen(j interface{}) bool {
	return i.second < j.(IntPair).second
}

func (i IntPair) EqualsTo(j interface{}) bool {
	return i.second == j.(IntPair).second
}

func (i IntPair) MinusInf() fib_heap.Value {
	return IntPair{-1e9, -Inf}
}

func IntPairComparator(a, b interface{}) int {
	if a.(IntPair).second == b.(IntPair).second {
		return 0
	}
	if a.(IntPair).second < b.(IntPair).second {
		return -1
	}
	return 1
}

type Edge struct {
	From Node
	To   Node
	Cost int64
}

const Inf = 1e18

type Graph struct {
	Size  int
	edges map[Node][]IntPair
}

func (g *Graph) AddEdge(e Edge) {
	if g.Size == 0 || g.edges == nil {
		g.init()
	}
	if g.edges[e.From] == nil {
		g.edges[e.From] = make([]IntPair, 0)
	}
	if g.edges[e.To] == nil {
		g.edges[e.To] = make([]IntPair, 0)
	}
	g.edges[e.From] = append(g.edges[e.From], IntPair{e.To, e.Cost})
	g.edges[e.To] = append(g.edges[e.To], IntPair{e.From, e.Cost})
}

func (g *Graph) init() {
	g.edges = make(map[Node][]IntPair)
}
