# fib_heap
Golang fibonacci heap + Dijkstra SSSP algorithm implementation

## Summary
Here is the implmentation for [Fibonacci Heap](https://en.wikipedia.org/wiki/Fibonacci_heap) on Go.

Also there is Dijkstra implementation on [Red-Black Tree](https://en.wikipedia.org/wiki/Red%E2%80%93black_tree), but the red-black tree set isn't mine

## Testing
There are unit tests for fib heap included, A large test for Dijkstra and I also submitted SSSP solutions based on this algo to [CF](https://codeforces.com/contest/20/submission/125209707)


## Details
  Here is time complexity of my [implementation](pkg/fib_heap/FibonacciHeap.go)
	
    Insertion     O(1)*
	GetMin        O(1)
	Merge         O(n) - details below
	ExtractMin    O(log n)
    Decrease      O(1)*
	DeleteKey     O(log n)
	  
    * - amortized
	Merge is O(n) because I use hash map to search for node in O(1)*
	If needed, can be replaced with other, tree-like structure, then
	Merge will be O(log n) but Decrease will be O(log n) too, and constant will increase
	It is usually more preferable to find decrease and delete keys from heap quickly then to merge
