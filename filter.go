package pipe

import (
	"sync"
)

// Filter takes message and forwards it if filter function return positive.
// If input channel is closed then output channel is closed.
// Creates a new channel with the same capacity as input.
//
// # Strategies
//
//   - Processing: Parallel
//   - Closing: Single
//   - Capacity: Same
//
// # Usages
//
//	// input := make(chan int, 4) with random values [1, 2, 3, 4]
//
//	output := Filter(func(value int) bool {
//		fmt.Print(value)
//	    return value % 2 == 0
//	}, input)
//
//	// stdout: 4 1 2 3
//	// output: [4 2]
func Filter[T any](filter func(T) bool, in <-chan T) <-chan T {
	out := make(chan T, cap(in))
	wg := sync.WaitGroup{}

	go func() {
		for {
			if in, ok := <-in; ok {
				wg.Add(1)
				go func() {
					if filter(in) {
						out <- in
					}
					wg.Done()
				}()
			} else {
				wg.Wait()
				close(out)
				break
			}
		}
	}()

	return out
}

// FilterSync takes message and forwards it if filter function return positive.
// If input channel is closed then output channel is closed.
// Creates a new channel with the same capacity as input.
//
// # Strategies
//
//   - Processing: Sync
//   - Closing: Single
//   - Capacity: Same
//
// # Usages
//
//	// input := make(chan int, 4) with random values [1, 2, 3, 4]
//
//	output := FilterSync(func(value int) bool {
//		fmt.Print(value)
//	    return value % 2 == 0
//	}, input)
//
//	// stdout: 4 1 2 3
//	// output: [2 4]
func FilterSync[T any](filter func(T) bool, in <-chan T) <-chan T {
	out := make(chan T, cap(in))
	queue := make(chan func() <-chan T, cap(in))
	wg := sync.WaitGroup{}

	go func() {
		for {
			if in, ok := <-in; ok {
				wg.Add(1)
				queue <- func() <-chan T {
					out := make(chan T)
					go func() {
						if filter(in) {
							out <- in
						}
						close(out)
						wg.Done()
					}()
					return out
				}
			} else {
				wg.Wait()
				close(queue)
				break
			}
		}
	}()

	go func() {
		for {
			if handler, ok := <-queue; ok {
				if data, ok := <-handler(); ok {
					out <- data
				}
			} else {
				close(out)
				break
			}
		}
	}()

	return out
}

// FilterSequential takes message and forwards it if filter function return positive.
// If input channel is closed then output channel is closed.
// Creates a new channel with the same capacity as input.
//
// # Strategies
//
//   - Processing: Sequential
//   - Closing: Single
//   - Capacity: Same
//
// # Usages
//
//	// input := make(chan int, 4) with random values [1, 2, 3, 4]
//
//	output := FilterSequential(func(value int) bool {
//		fmt.Print(value)
//	    return value % 2 == 0
//	}, input)
//
//	// stdout: 1 2 3 4
//	// output: [2 4]
func FilterSequential[T any](filter func(T) bool, in <-chan T) <-chan T {
	out := make(chan T, cap(in))

	go func() {
		for {
			if in, ok := <-in; ok {
				if filter(in) {
					out <- in
				}
			} else {
				close(out)
				break
			}
		}
	}()

	return out
}
