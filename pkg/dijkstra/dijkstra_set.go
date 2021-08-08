package dijkstra

import (
	treeset "github.com/emirpasic/gods/sets/treeset"
)

func DijkstraSet(g Graph, start int) (dist []int64, p []Node) {
	dist = make([]int64, g.Size+1)
	p = make([]int, g.Size+1)
	for i, _ := range dist {
		dist[i] = Inf
	}
	dist[start] = 0
	var set *treeset.Set = treeset.NewWith(IntPairComparator)
	set.Add(IntPair{start, 0})

	for !set.Empty() {
		it := set.Iterator()
		it.First()
		cur := it.Value().(IntPair)
		v := cur.first
		set.Remove(cur)

		for _, pair := range g.edges[v] {
			to := pair.first
			cost := pair.second
			if dist[to] > dist[v]+cost {
				set.Remove(IntPair{to, dist[to]})
				dist[to] = dist[v] + cost
				p[to] = v
				set.Add(IntPair{to, dist[to]})
			}

		}
	}

	return
}

func ShortestPathSet(g Graph, start int, finish int) (dist int64, path []Node) {
	d, pr := DijkstraSet(g, start)
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
