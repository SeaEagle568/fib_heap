package fib_heap

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"
)

type pair struct {
	first  int
	second int
}

func (i pair) EqualsTo(j interface{}) bool {
	return i.first == j.(pair).first
}
func (i pair) LessThen(j interface{}) bool {
	if i.first != j.(pair).first {
		return i.first < j.(pair).first
	} else {
		return i.second < j.(pair).second
	}
}

func (i pair) MinusInf() Value {
	return pair{-1e9, -1e9}
}

func EmptyTest(t *testing.T, pq *FibonacciHeap) {
	ass := assert.New(t)
	ass.Equal(pq.IsEmpty(), true, "fib heap not empty")
	ass.Equal(pq.GetSize(), 0, "fib heap size != 0")

	val, err := pq.FindMin()
	if err == nil {
		t.Log(val)
		t.Error("Empty list extraxtion")
	}
	val, err = pq.ExtractMin()
	if err == nil {
		t.Log(val)
		t.Error("Empty list extraction")
	}
}

func OneElement(t *testing.T, pq *FibonacciHeap, elem pair) {
	ass := assert.New(t)
	pq.Insert(elem)
	ass.Equal(pq.IsEmpty(), false, "fib heap empty after insertion")
	ass.Equal(pq.GetSize(), 1, "fib heap size != 1 after insertion")
	val, err := pq.FindMin()
	if err != nil {
		t.Error("error finding min")
	}
	ass.Equal(val, elem, "FindMin failed. Expected {}, found {}", elem, val)
	val, err = pq.ExtractMin()
	if err != nil {
		t.Error("Error extracting min")
	}
	ass.Equal(val, elem, "FindMin failed. Expected {}, found {}", elem, val)
	EmptyTest(t, pq)
}

func TwoElements(t *testing.T, pq *FibonacciHeap, elem1 pair, elem2 pair) {
	ass := assert.New(t)
	min, max := elem1, elem2
	if elem2.LessThen(elem1) {
		min, max = max, min
	}
	pq.Insert(max)
	pq.Insert(min)
	ass.Equal(pq.IsEmpty(), false, "fib heap empty after insertion")
	ass.Equal(pq.GetSize(), 2, "fib heap size != 2 after insertion")

	val, err := pq.FindMin()
	if err != nil {
		t.Error("error finding min")
	}
	ass.Equal(val, min, "FindMin failed. Expected {}, found {}", min, val)
	val, err = pq.ExtractMin()
	if err != nil {
		t.Error("Error extracting min")
	}
	ass.Equal(val, min, "FindMin failed. Expected {}, found {}", min, val)

	val, err = pq.FindMin()
	if err != nil {
		t.Error("error finding min")
	}
	ass.Equal(val, max, "FindMin failed. Expected {}, found {}", max, val)
	val, err = pq.ExtractMin()
	if err != nil {
		t.Error("Error extracting min")
	}
	ass.Equal(val, max, "FindMin failed. Expected {}, found {}", max, val)
}

func RandomTest(t *testing.T, heap *FibonacciHeap, n int, arr []pair) {
	ass := assert.New(t)
	ass.Equal(heap.GetSize(), n, "wrong fibheap size")
	for _, elem := range arr {
		val, err := heap.ExtractMin()
		if err != nil {
			t.Error("error extrcting min")
		}
		ass.Equal(val.(pair).first, elem.first)
	}
	EmptyTest(t, heap)
}

func GenerateCoolHeap(from int, to int) (heap *FibonacciHeap, n int, arr []pair) {
	if to < from {
		to, from = from, to
	}
	n = to - from
	heap = &FibonacciHeap{}
	rand.Seed(1488228)
	arr = make([]pair, n)
	for i := 0; i < n; i++ {
		arr[i] = pair{from + i, rand.Int()}
	}
	shuffled := make([]pair, n)
	copy(shuffled, arr)
	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
	for _, val := range shuffled {
		heap.Insert(val)
	}
	return
}
func NRandomTest(t *testing.T, n int) {
	for i := 0; i < n; i++ {
		randSize := rand.Intn(5000) + 100
		heap, sz, arr := GenerateCoolHeap(0, randSize)
		RandomTest(t, heap, sz, arr)
	}
}

func NMergeTest(t *testing.T, n int) {
	for i := 0; i < n; i++ {
		randSize := rand.Intn(5000) + 100
		heap, sz, arr := GenerateCoolHeap(0, randSize)
		heap2, sz2, arr2 := GenerateCoolHeap(randSize+1, randSize+1+rand.Intn(5000)+100)
		heap.Merge(heap2)
		sz += sz2
		arr = append(arr, arr2...)
		RandomTest(t, heap, sz, arr)
	}
}

func remove(s []pair, i int) []pair {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
func NDeleteTest(t *testing.T, n int) {
	for i := 0; i < n; i++ {
		randSize := rand.Intn(5000) + 100
		//randSize := 10
		heap, sz, arr := GenerateCoolHeap(0, randSize)
		for j := 0; j < 50; j++ {
			pos := rand.Intn(len(arr))
			err := heap.DeleteKey(arr[pos])
			if err != nil {
				t.Error(err)
			}
			arr = remove(arr, pos)
			sz--
		}
		sort.Slice(arr, func(i, j int) bool {
			return arr[i].first < arr[j].first
		})
		RandomTest(t, heap, sz, arr)
	}
}

func NChangeTest(t *testing.T, n int) {
	for i := 0; i < n; i++ {
		randSize := rand.Intn(5000) + 100
		heap, sz, arr := GenerateCoolHeap(0, randSize)
		for j := 0; j < 100; j++ {
			pos := rand.Intn(sz)
			var delta = 0
			if arr[pos].first != 0 {
				delta = rand.Intn(arr[pos].first)
			}
			value, second := arr[pos].first-delta, rand.Int()
			err := heap.Decrease(arr[pos], pair{value, second})
			if err != nil {
				t.Error(err)
			}
			arr[pos] = pair{value, second}
		}
		sort.Slice(arr, func(i, j int) bool {
			return arr[i].first < arr[j].first
		})
		RandomTest(t, heap, sz, arr)
	}
}

func Test(t *testing.T) {
	EmptyTest(t, &FibonacciHeap{})
	OneElement(t, &FibonacciHeap{}, pair{1, 1})
	TwoElements(t, &FibonacciHeap{}, pair{1, 100}, pair{100, 1})

	NRandomTest(t, 20)
	NMergeTest(t, 20)
	NChangeTest(t, 20)
	NDeleteTest(t, 20)
}
