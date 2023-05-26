package pipe

// Wait waits for the input channel to close and sends a signal to the returned channel.
//
// # Example
//
//	<-Wait(input1)
//	select {
//		case <-Wait(input2):
//		case <-Wait(input3):
//	}
//	// Will executed after input1 closed and input2 or input3 closed
func Wait[T any](in <-chan T) <-chan struct{} {
	q := make(chan struct{})
	go func() {
		for {
			if _, ok := <-in; !ok {
				break
			}
		}
		q <- struct{}{}
		close(q)
	}()
	return q
}

// WaitAll waits for all input channels to close and sends a signal to the returned channel.
//
// # Example
//
//	<-WaitAll(input1, input2)
//	// Will executed after input1 AND input2 closed
//
//	// It's equal:
//	<-Wait(input1)
//	<-Wait(input2)
func WaitAll[T any](in ...<-chan T) <-chan struct{} {
	q := make(chan struct{})
	go func() {
		for i := range in {
			for {
				if _, ok := <-in[i]; !ok {
					break
				}
			}
		}
		q <- struct{}{}
		close(q)
	}()
	return q
}

// WaitAny waits for one of the input channels to close and sends a signal to the returned channel.
// All other channels are read to the end in the background.
//
// # Example
//
//	<-WaitAny(input1, input2)
//	// Will executed after input1 OR input2 closed
//
//	// It's equal:
//	select {
//		case <-Wait(input1):
//		case <-Wait(input2):
//	}
func WaitAny[T any](in ...<-chan T) <-chan struct{} {
	q := make(chan struct{})
	go func() {
		queue := make(chan struct{}, len(in))
		for i := range in {
			i := i
			go func() {
				for {
					if _, ok := <-in[i]; !ok {
						queue <- struct{}{}
						break
					}
				}
			}()
		}
		<-queue
		q <- struct{}{}
	}()
	return q
}
