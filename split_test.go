package pipe

import (
	"errors"
	"sync"
	"testing"

	"github.com/msacore/pipe/test"
)

func TestSplit(t *testing.T) {
	t.Run("Parallel", func(t *testing.T) {
		test.Suit(t, func(epipe chan error) ([]<-chan int, chan error) {
			pipe := test.Generator(0, 64, 16)

			pipes := Split(2, pipe)
			pipes[0], epipe = test.AssertCount("count 1", pipes[0], epipe, 64)
			pipes[1], epipe = test.AssertCount("count 2", pipes[1], epipe, 64)

			return pipes, epipe
		})
	})

	t.Run("Sync", func(t *testing.T) {
		const channelCapacity = 64
		const dataCount = channelCapacity * 8
		const outputsCount = 16

		input := make(chan int, channelCapacity)
		go func() {
			for i := 0; i < dataCount; i++ {
				input <- i
			}
			close(input)
		}()

		outs := SplitSync(outputsCount, input)
		errs := make([]error, outputsCount)

		wg := sync.WaitGroup{}
		for i := 0; i < outputsCount; i++ {
			i := i
			wg.Add(1)
			go func() {
				max := -1
				for val := range outs[i] {
					if val > max {
						max = val
					} else {
						errs[i] = errors.New("the data order is broken")
					}
				}
				if max < dataCount-1 {
					errs[i] = errors.New("not all data has been processed")
				}
				wg.Done()
			}()
		}
		wg.Wait()

		for _, err := range errs {
			if err != nil {
				t.Fatal(err)
			}
		}
	})

	t.Run("Sync", func(t *testing.T) {
		test.Suit(t, func(epipe chan error) ([]<-chan int, chan error) {
			pipe := test.Generator(0, 64, 16)

			pipes := SplitSync(2, pipe)
			pipes[0], epipe = test.AssertCount("count 1", pipes[0], epipe, 64)
			pipes[0], epipe = test.AssertOrderAsc("ordering 1", pipes[0], epipe)
			pipes[1], epipe = test.AssertCount("count 2", pipes[1], epipe, 64)
			pipes[1], epipe = test.AssertOrderAsc("ordering 2", pipes[1], epipe)

			return pipes, epipe
		})
	})

	t.Run("Sequential", func(t *testing.T) {
		test.Suit(t, func(epipe chan error) ([]<-chan int, chan error) {
			pipe := test.Generator(0, 64, 16)

			pipes := SplitSync(2, pipe)
			pipes[0], epipe = test.AssertCount("count 1", pipes[0], epipe, 64)
			pipes[0], epipe = test.AssertOrderAsc("ordering 1", pipes[0], epipe)
			pipes[1], epipe = test.AssertCount("count 2", pipes[1], epipe, 64)
			pipes[1], epipe = test.AssertOrderAsc("ordering 2", pipes[1], epipe)

			return pipes, epipe
		})
	})
}
