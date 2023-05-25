package pipe

import "sync"

// Take message and convert it into another type by map function.
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
//	// input := make(chan int, 4) with random values [1, 2, 3]
//	output := Map(input, func(value int) string {
//	    fmt.Print(value)
//	    return fmt.Sprintf("val: %d", value)
//	})
//	// stdout: 2 1 3
//	// output: ["val: 2", "val: 1", "val: 3"]
func Map[Tin, Tout any](in <-chan Tin, mapper func(Tin) Tout) <-chan Tout {
	out := make(chan Tout, cap(in))
	wg := sync.WaitGroup{}

	go func() {
		for {
			if in, ok := <-in; ok {
				wg.Add(1)
				go func() {
					out <- mapper(in)
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

// Take message and convert it into another type by map function.
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
//	// input := make(chan int, 4) with random values [1, 2, 3]
//	output := Map(input, func(value int) string {
//	    fmt.Print(value)
//	    return fmt.Sprintf("val: %d", value)
//	})
//	// stdout: 2 1 3
//	// output: ["val: 1", "val: 2", "val: 3"]
func MapSync[Tin, Tout any](in <-chan Tin, mapper func(Tin) Tout) <-chan Tout {
	out := make(chan Tout, cap(in))
	queue := make(chan func() <-chan Tout, cap(in))
	wg := sync.WaitGroup{}

	go func() {
		for {
			if in, ok := <-in; ok {
				wg.Add(1)
				queue <- func() <-chan Tout {
					out := make(chan Tout)
					go func() {
						out <- mapper(in)
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
				out <- <-handler()
			} else {
				close(out)
				break
			}
		}
	}()

	return out
}

// Take message and convert it into another type by map function.
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
//	// input := make(chan int, 4) with random values [1, 2, 3]
//	output := Map(input, func(value int) string {
//	    fmt.Print(value)
//	    return fmt.Sprintf("val: %d", value)
//	})
//	// stdout: 1 2 3
//	// output: ["val: 1", "val: 2", "val: 3"]
func MapSequential[Tin, Tout any](in <-chan Tin, mapper func(Tin) Tout) <-chan Tout {
	out := make(chan Tout, cap(in))

	go func() {
		for {
			if in, ok := <-in; ok {
				out <- mapper(in)
			} else {
				close(out)
				break
			}
		}
	}()

	return out
}