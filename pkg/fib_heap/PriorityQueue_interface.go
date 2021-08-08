package fib_heap

type PriorityQueue interface {
	FindMin() (Value, error)
	Insert(x Value)
	ExtractMin() (Value, error)
	Decrease(x Value, newVal Value) error
	IsEmpty() bool
	DeleteKey(x Value) error
}
