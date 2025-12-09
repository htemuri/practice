/*
A Fixed-Size Ring Buffer (no generics)

Skills: slices, modular arithmetic, struct design, error handling.
Why: Ring buffers show up everywhere in low-latency systems and telemetry pipelines.

Requirements:

	Push(item) overwrites oldest when full
	Pop() returns oldest
	Internal storage is a slice
	Track head/tail indices
	100–200 lines max
*/
package main

import "fmt"

type RingBuffer[T any] struct {
	buffer []T
	tail   int
	head   int
	empty  bool
}

func NewRingBuffer[T any](size int) *RingBuffer[T] {
	return &RingBuffer[T]{
		buffer: make([]T, size),
		empty:  true,
		head:   0,
		tail:   -1,
	}
}

func (ring *RingBuffer[T]) push(item T) {
	ringSize := len(ring.buffer)
	ring.tail = (ring.tail + 1) % ringSize

	if ring.tail == ring.head {
		ring.head = (ring.head + 1) % ringSize
	}

	ring.buffer[ring.tail] = item
	ring.empty = false
}

func (ring *RingBuffer[T]) pop() (T, error) {
	// if empty
	var zero T
	if ring.empty {
		return zero, fmt.Errorf("empty ring")
	}

	// if want to return head and increment
	if ring.head == ring.tail {
		ring.empty = true
		return ring.buffer[ring.head], nil
	}

	popped := ring.buffer[ring.head]
	ring.head = (ring.head + 1) % len(ring.buffer)
	return popped, nil
}

func main() {
	buffer := NewRingBuffer[int](5)
	buffer.push(3)
	fmt.Println(buffer)
	buffer.push(1)
	fmt.Println(buffer)
	buffer.push(234)
	fmt.Println(buffer)
	buffer.push(44)
	fmt.Println(buffer)
	fmt.Println(buffer.pop())
	fmt.Println(buffer.pop())
	fmt.Println(buffer.pop())
	fmt.Println(buffer.pop())
	fmt.Println(buffer.pop())
}
