/**
This is Dijkstra on Fibonacci Heap realisation in one file
that solves this Codeforces Problem https://codeforces.com/contest/20/problem/C
(It gives TL but only because GO is a little slow without optimizations)

@author SeaEagle
*/

package main

import (
	"errors"
	"fmt"
)

func main() {
	var n, m int
	fmt.Scan(&n, &m)
	g := Graph{Size: n}
	for i := 0; i < m; i++ {
		var from, to, cost int
		fmt.Scan(&from, &to, &cost)
		g.AddEdge(Edge{From: from, To: to, Cost: int64(cost)})
	}
	dist, path := ShortestPathHeap(&g, 1, n)
	if dist == Inf {
		fmt.Println(-1)
	} else {
		for _, val := range path {
			fmt.Printf("%v ", val)
		}
	}
}

type GNode = int

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

func (i IntPair) MinusInf() Value {
	return IntPair{-1e9, -Inf}
}

type Edge struct {
	From GNode
	To   GNode
	Cost int64
}

const Inf = 1e18

type Graph struct {
	Size  int
	edges map[GNode][]IntPair
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
	g.edges = make(map[GNode][]IntPair)
}

func DijkstraHeap(g *Graph, start int) (dist []int64, p []GNode) {
	dist = make([]int64, g.Size+1)
	p = make([]int, g.Size+1)
	for i := range dist {
		dist[i] = Inf
	}
	dist[start] = 0
	heap := FibonacciHeap{}
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

func ShortestPathHeap(g *Graph, start int, finish int) (dist int64, path []GNode) {
	d, pr := DijkstraHeap(g, start)
	dist = d[finish]
	if dist == Inf {
		return
	}
	path = make([]GNode, 0)
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

//Value is an interface every heap value must implement
//Three self-explanatory methods
type Value interface {
	LessThen(j interface{}) bool
	EqualsTo(j interface{}) bool
	MinusInf() Value
}

//Node is a struct for node, all Values need to implement Value interface
type Node struct {
	key    Value
	parent *Node
	child  *Node //Only one child, we can iterate children cyclical using left right pointers
	left   *Node //Cyclic doubly-linked-list
	right  *Node
	degree int  //Number of child ren
	marked bool //If we already cut this node's child
}

//FibonacciHeap is represented with pointer to minimum, size
//and hashmap (value, pointer to node with that value)
//THAT MEANS TO USE DECREASE CORRECTLY YOU NEED UNIQUE ENTIES
//It could be done with some sort of hashes in your structure
type FibonacciHeap struct {
	min   *Node
	size  int
	nodes map[Value]*Node
}

//Decrease method accepts old and new value and changes that in heap
//If there is no entry that matches OR new value is bigger than old
//A not-nil error is returned
func (f *FibonacciHeap) Decrease(x Value, newVal Value) error {
	if f.nodes[x] == nil {
		return errors.New("no such element in heap")
	}
	if !newVal.LessThen(x) && !newVal.EqualsTo(x) {
		return errors.New("cannot increase element")
	}
	node := f.nodes[x]
	delete(f.nodes, x) //Delete entry from map and put updated one
	f.nodes[newVal] = node
	node.key = newVal

	if node.parent == nil || node.parent.key.LessThen(newVal) { //If it is root or we don't break the heap
		f.updateMin(node)
		return nil //Update minimum and return
	}
	parent := node.parent
	f.cut(node)            //Cut the subtree out and make it a new root
	f.cascadingCut(parent) //Do it with all parents if needed
	return nil
}

//cascadingCut is a method that goes up a heap,
//and cuts every parent who is <b>marked</b>
func (f *FibonacciHeap) cascadingCut(node *Node) {
	for node.parent != nil && node.marked {
		parent := node.parent
		f.cut(node)
		node = parent
	}
	if node.parent != nil {
		node.marked = true
		f.updateMin(node.parent)
	}
	f.updateMin(node)
}

//cut is a method to cut out a sub-heap from parent
//and take it out to roots level
func (f *FibonacciHeap) cut(node *Node) {
	node.left.right, node.right.left = node.right, node.left //cut the node from linked list
	node.parent.degree--
	if node.parent.child == node { //If this node is the main child, we need to make a new main child
		if node.parent.degree == 0 {
			node.parent.child = nil
		} else {
			node.parent.child = node.right
		}
	}
	node.left, node.right = node, node //inserting node into top level linked list
	node.parent = nil
	f.uniteLists(f.min, node)
	f.updateMin(node)
}

//DeleteKey accepts a Value to delete
//It is very simple - just Decrease the needed Value
//And ExtractMin-s it
func (f *FibonacciHeap) DeleteKey(x Value) error {
	err := f.Decrease(x, x.MinusInf())
	if err != nil {
		return err
	}
	_, err = f.ExtractMin()
	return err
}

//FindMin returns the minimum value of the heap
//Just returning the pointer to min
func (f *FibonacciHeap) FindMin() (Value, error) {
	if !f.IsEmpty() {
		return f.min.key, nil
	} else {
		return nil, errors.New("unable to get minimum from an empty heap")
	}
}

//Insert inserts certain Value into heap
//By just inserting it into root-level linked-list
//To the right from minimum
func (f *FibonacciHeap) Insert(x Value) {
	if f.IsEmpty() {
		f.init()
		f.min.key = x
		f.min.left, f.min.right = f.min, f.min
		f.nodes[x] = f.min
	} else {
		newNode := &Node{x, nil, nil, f.min, f.min.right, 0, false} //Create a new node
		f.min.right.left = newNode
		f.min.right = newNode
		f.updateMin(newNode)
		f.nodes[x] = newNode
	}
	f.size++
}

//ExtractMin deletes the minimum value of the heap and returns it
//If heap is empty there comes an error
//Implementation is not that simple, we cut all children
//and make them new roots, after that we need to rebuild out heap
//with consolidate method
func (f *FibonacciHeap) ExtractMin() (Value, error) {
	if f.IsEmpty() {
		return nil, errors.New("unable to get minimum from an empty heap")
	}
	res := f.min.key
	delete(f.nodes, f.min.key)
	if f.GetSize() == 1 { //If it was the only node return
		f.min = nil
		f.size = 0
		return res, nil
	}
	f.cutChildren(f.min, f.min.child)
	f.min.left.right = f.min.right
	f.min.right.left = f.min.left
	f.min = f.min.right //Move min pointer to the right, consolidation will fin new min
	f.consolidate()
	f.size--
	return res, nil
}

//cutChildren is a fun-named method that cuts children from parent
//It iterates throw all children to invalidate parent pointer
//And then just makes them new roots
func (f *FibonacciHeap) cutChildren(father *Node, child *Node) {
	if child == nil {
		return
	}
	start := child
	start.parent = nil
	cur := child.right
	for cur != start {
		cur.parent = nil
		cur = cur.right
	}
	father.child = nil
	f.uniteLists(father, child)
}

//consolidate is a heavy O(log n) function
//What it does is just find two root trees with same degree
//And then hang the bigger (by node.key) one to smaller
func (f *FibonacciHeap) consolidate() {
	var used = make([]*Node, f.size) //We use a slice to track collisions in node.degree
	used[f.min.degree] = f.min
	cur := f.min.right

	for used[cur.degree] != cur { //We always go right, so if we placed something to slice, made a lap, and nothing changed, consolidation is finished
		f.updateMin(cur)
		if used[cur.degree] == nil { //If yet no other node with same degree recorder - record current
			used[cur.degree] = cur
			cur = cur.right
		} else {

			busy := used[cur.degree]
			father, son := cur, busy
			if busy.key.LessThen(cur.key) { //make father point to lighter node, son to heavier one
				father, son = son, father
			} else if father.key.EqualsTo(son.key) { //make sure f.min is always father
				if son == f.min {
					father, son = son, father
				}
			}
			son.left.right = son.right
			son.right.left = son.left //cut the son from his local linked-list

			next := cur.right //remember next to be right from current cur, it can change later

			if father.child == nil { //If father has no children - son is the first
				father.child = son
				son.left, son.right = son, son
			} else { //else integrate son into children linked-list
				son.left, son.right = father.child, father.child.right
				father.child.right.left = son
				father.child.right = son
			}

			used[cur.degree] = nil
			son.parent = father
			father.degree++
			cur = next
		}
		f.updateMin(cur)
	}

}
func (f *FibonacciHeap) updateMin(comp *Node) {
	if comp.key.LessThen(f.min.key) {
		f.min = comp
		return
	}
	if comp.key.EqualsTo(f.min.key) {
		if comp.parent == nil && f.min.parent != nil {
			f.min = comp
			return
		}
		for f.min.parent != nil {
			f.min = f.min.parent
		}
	}

}

//Merge just merges two fibonacci heaps together
//It could be O(1) if we dont need to hold relevant hashmap
//That helps to do search in O(1)*
func (f *FibonacciHeap) Merge(heap *FibonacciHeap) {
	if heap.size == 0 {
		return
	}
	if f.size == 0 { //if our heap is zero-sized, just change pointers to another list
		f.min = heap.min
		f.size = heap.size
	} else {
		f.uniteLists(f.min, heap.min) //Unites two lists
		f.size += heap.size
		f.updateMin(heap.min)
	}
	for k, v := range heap.nodes { //O(n) code here!!
		f.nodes[k] = v
	}
}

//uniteLists is simple function that unites two cyclic doubly-linked-lists into one
func (f *FibonacciHeap) uniteLists(first *Node, second *Node) {
	if second == nil || first == nil {
		return
	}
	first.left.right = second.right
	second.right.left = first.left
	first.left = second
	second.right = first
}

//GetSize returns current size of the heap (number of nodes)
func (f *FibonacciHeap) GetSize() int {
	return f.size
}

//IsEmpty returns true if fibheap is empty
func (f *FibonacciHeap) IsEmpty() bool {
	return f.size == 0
}

//init function creates an empty node and initialize hashmap
func (f *FibonacciHeap) init() {
	f.min = &(Node{})
	f.nodes = make(map[Value]*Node)
	f.size = 0
}
