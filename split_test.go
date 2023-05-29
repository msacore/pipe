package pipe

import (
	"errors"
	"sync"
	"testing"
)

func TestSplit(t *testing.T) {
	t.Run("Parallel", func(t *testing.T) {
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

		outs := Split(outputsCount, input)

		checks := make([][]bool, outputsCount)
		for i := 0; i < outputsCount; i++ {
			checks[i] = make([]bool, dataCount)
		}

		wg := sync.WaitGroup{}
		for i := 0; i < outputsCount; i++ {
			i := i
			wg.Add(1)
			go func() {
				for val := range outs[i] {
					checks[i][int(val)] = true
				}
				wg.Done()
			}()
		}
		wg.Wait()

		for i := 0; i < outputsCount; i++ {
			i := i
			for _, check := range checks[i] {
				if !check {
					t.Fatal("not all data has been processed")
				}
			}
		}
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

	t.Run("Sequential", func(t *testing.T) {
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

		outs := SplitSequential(outputsCount, input)
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
}
