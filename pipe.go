// Package pipe provides a convenient way to work with Go channels and simple construction of pipelines with a wide
// range of logic gates.
//
// # Processing Strategies
//
// Some functions have different channel processing algorithms. To ensure maximum performance, it is recommended
// to use the original function. However, specific algorithms can help in cases where you are faced with a race of
// threads or you need to output data strictly in the same order in which you received them.
//
//   - Parallel - Each handler is executed in its own goroutine and there is no guarantee that the output order
//     will be consistent. Recommended for best performance.
//   - Sync - Each handler executes in its own goroutine, but the result of the youngest goroutine waits for the oldest
//     goroutine to finish before being passed to the output stream. To prevent memory leaks, the strategy will wait if
//     there is more waiting data than the capacity of the output channel. Recommended if you want to get the output
//     data in the same order as the input data.
//   - Sequential - Each handler is executed sequentially, one after the other. Keeps the order of the output data
//     equal to the order of the input data. Recommended if it is necessary to exclude the race of threads between
//     handlers.
//
// If the input channel capacity is 0 (no bandwidth), then any strategy will act as Sequential behavior.
//
// # Closing Strategies
//
// Each function has one of several strategies for closing output channels. Understanding will help you understand
// how and when your pipeline closes.
//
//   - Single - Suitable only for functions with one input. If the input channel is closed, then the output channels
//     are closed.
//   - All - If all input channels are closed, then the output channels are closed.
//   - Any - If one of the input channels is closed, the output channels are closed. All other channels will be read
//     to the end in the background.
//
// # Capacity Strategies
//
// Each function creates new output channels with the capacity corresponding to a specific strategy.
//
//   - Same - Suitable only for functions with one input channel. The output channels will have a capacity equal to
//     the input channel.
//   - Mult - Suitable only for functions with one input channel. The output channels will have a capacity equal to
//     the input channel multiplied by N.
//   - Min - The output channels will have a capacity equal to the minimum capacity of the input channels.
//   - Max - The output channels will have a capacity equal to the maximum capacity of the input channels.
//   - Sum - The output channels will have a capacity equal to the sum of capacities of the input channels.
package pipe

import "sync"

// ========== DEV ========== //

func Filter[T any](in <-chan T, filter func(T) bool) <-chan T {
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

func FilterSequential[T any](in <-chan T, filter func(T) bool) <-chan T {
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

// Watch calls callback for every item in sequence without ordering.
func Watch[T any](in <-chan T, callback func(T)) <-chan T {
	out := make(chan T, cap(in))
	wg := sync.WaitGroup{}
	go func() {
		for {
			if in, ok := <-in; ok {
				wg.Add(1)
				go func() {
					callback(in)
					out <- in
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

// WatchSequential calls callback for every item in sequence with input order.
func WatchSequential[T any](in <-chan T, callback func(T)) <-chan T {
	out := make(chan T, cap(in))
	go func() {
		for {
			if in, ok := <-in; ok {
				callback(in)
				out <- in
			} else {
				close(out)
				break
			}
		}
	}()
	return out
}

// Wait is terminate function. It will block current goroutine until input channel is closed.
func Wait[T any](in <-chan T) {
	for {
		if _, ok := <-in; !ok {
			break
		}
	}
}

// func Generate[T any](generator func() (T, bool)) <-chan T {
// }
