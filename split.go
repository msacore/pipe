package pipe

import "sync"

// Split takes a number of output channels and input channel, and forwards the input
// messages to all output channels.
// There is no guarantee that the message will be sent to the output channels in the
// sequence in which they are provided.
// If input channel is closed then all output channels are closed.
// Creates new channels with the same capacity as input.
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
//
//	outs := Split(2, input)
//
//	// The gaps demonstrate uneven recording in the channels
//	// outs[0]: [2,    1, 3   ]
//	// outs[1]: [   1, 3,    2]
func Split[T any](n int, in <-chan T) []<-chan T {
	outs := make([]chan T, n)
	for i := 0; i < n; i++ {
		outs[i] = make(chan T, cap(in))
	}
	wg := sync.WaitGroup{}

	go func() {
		for {
			if in, ok := <-in; ok {
				for i := 0; i < n; i++ {
					i := i
					wg.Add(1)
					go func() {
						outs[i] <- in
						wg.Done()
					}()
				}
			} else {
				wg.Wait()
				for i := 0; i < n; i++ {
					close(outs[i])
				}
				break
			}
		}
	}()

	outsR := make([]<-chan T, n)
	for i := 0; i < n; i++ {
		outsR[i] = outs[i]
	}
	return outsR
}

// Split2 - alias for [Split]
func Split2[T any](in <-chan T) (out1, out2 <-chan T) {
	outs := Split(2, in)
	return outs[0], outs[1]
}

// Split3 - alias for [Split]
func Split3[T any](in <-chan T) (out1, out2, out3 <-chan T) {
	outs := Split(3, in)
	return outs[0], outs[1], outs[2]
}

// SplitSync takes a number of output channels and input channel, and forwards the input
// messages to all output channels.
// There is no guarantee that the message will be sent to the output channels in the
// sequence in which they are provided.
// If input channel is closed then all output channels are closed.
// Creates new channels with the same capacity as input.
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
//
//	outs := SplitSync(2, input)
//
//	// The gaps demonstrate uneven recording in the channels
//	// outs[0]: [1,    2, 3   ]
//	// outs[1]: [   1, 2,    3]
func SplitSync[T any](n int, in <-chan T) []<-chan T {
	outs := make([]chan T, n)
	queues := make([]chan func() <-chan T, cap(in))
	for i := 0; i < n; i++ {
		outs[i] = make(chan T, cap(in))
	}
	wg := sync.WaitGroup{}

	go func() {
		for {
			if in, ok := <-in; ok {
				for i := 0; i < n; i++ {
					i := i
					wg.Add(1)
					queues[i] <- func() <-chan T {
						out := make(chan T)
						go func() {
							out <- in
							close(out)
							wg.Done()
						}()
						return out
					}
				}
			} else {
				wg.Wait()
				for i := 0; i < n; i++ {
					close(queues[i])
				}
				break
			}
		}
	}()

	for i := 0; i < n; i++ {
		i := i
		go func() {
			for {
				if handler, ok := <-queues[i]; ok {
					if data, ok := <-handler(); ok {
						outs[i] <- data
					}
				} else {
					close(outs[i])
					break
				}
			}
		}()
	}

	outsR := make([]<-chan T, n)
	for i := 0; i < n; i++ {
		outsR[i] = outs[i]
	}
	return outsR
}

// SplitSync2 - alias for [SplitSync]
func SplitSync2[T any](in <-chan T) (out1, out2 <-chan T) {
	outs := SplitSync(2, in)
	return outs[0], outs[1]
}

// SplitSync3 - alias for [SplitSync]
func SplitSync3[T any](in <-chan T) (out1, out2, out3 <-chan T) {
	outs := SplitSync(3, in)
	return outs[0], outs[1], outs[2]
}

// SplitSequential takes a number of output channels and input channel, and forwards the input
// messages to all output channels.
// The message will be sent to the output channels in the following sequence.
// If input channel is closed then all output channels are closed.
// Creates new channels with the same capacity as input.
//
// Be aware, if one of the output channels is blocked, then all other output channels will wait.
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
//
//	outs := SplitSequential(2, input)
//
//	// The gaps demonstrate uneven recording in the channels
//	// outs[0]: [1,    2,    3   ]
//	// outs[1]: [   1,    2,    3]
func SplitSequential[T any](n int, in <-chan T) []<-chan T {
	outs := make([]chan T, n)
	for i := 0; i < n; i++ {
		outs[i] = make(chan T, cap(in))
	}

	go func() {
		for {
			if in, ok := <-in; ok {
				for i := 0; i < n; i++ {
					outs[i] <- in
				}
			} else {
				for i := 0; i < n; i++ {
					close(outs[i])
				}
				break
			}
		}
	}()

	outsR := make([]<-chan T, n)
	for i := 0; i < n; i++ {
		outsR[i] = outs[i]
	}
	return outsR
}

// SplitSequential2 - alias for [SplitSequential]
func SplitSequential2[T any](in <-chan T) (out1, out2 <-chan T) {
	outs := SplitSequential(2, in)
	return outs[0], outs[1]
}

// SplitSequential3 - alias for [SplitSequential]
func SplitSequential3[T any](in <-chan T) (out1, out2, out3 <-chan T) {
	outs := SplitSequential(3, in)
	return outs[0], outs[1], outs[2]
}
