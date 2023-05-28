package pipe

import (
	"testing"
)

func TestFilter(t *testing.T) {
	t.Run("Parallel", func(t *testing.T) {
		const channelCapacity = 64
		const dataCount = channelCapacity * 8

		input := make(chan int, channelCapacity)
		go func() {
			for i := 0; i < dataCount; i++ {
				input <- i
			}
			close(input)
		}()

		output := Filter(func(val int) bool {
			return val%2 == 0
		}, input)

		count := 0
		for val := range output {
			count++
			if val%2 != 0 {
				t.Fatal("false positive filtration")
			}
		}

		if count != dataCount/2 {
			t.Fatal("not all data has been processed")
		}
	})

	t.Run("Sync", func(t *testing.T) {
		const channelCapacity = 64
		const dataCount = channelCapacity * 8

		input := make(chan int, channelCapacity)
		go func() {
			for i := 0; i < dataCount; i++ {
				input <- i
			}
			close(input)
		}()

		output := FilterSync(func(val int) bool {
			return val%2 == 0
		}, input)

		max := -1
		count := 0
		for val := range output {
			count++
			if val%2 != 0 {
				t.Fatal("false positive filtration")
			}
			if val > max {
				max = val
			} else {
				t.Fatal("the data order is broken")
			}
		}

		if count != dataCount/2 {
			t.Fatal("not all data has been processed")
		}
	})

	t.Run("Sequential", func(t *testing.T) {
		const channelCapacity = 64
		const dataCount = channelCapacity * 8

		input := make(chan int, channelCapacity)
		go func() {
			for i := 0; i < dataCount; i++ {
				input <- i
			}
			close(input)
		}()

		output := FilterSequential(func(val int) bool {
			return val%2 == 0
		}, input)

		max := -1
		count := 0
		for val := range output {
			count++
			if val%2 != 0 {
				t.Fatal("false positive filtration")
			}
			if val > max {
				max = val
			} else {
				t.Fatal("the data order is broken")
			}
		}

		if count != dataCount/2 {
			t.Fatal("not all data has been processed")
		}
	})
}
