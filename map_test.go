package pipe

import (
	"testing"
)

func TestMap(t *testing.T) {
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

		output := Map(func(val int) int32 {
			return int32(val)
		}, input)

		checks := make([]bool, dataCount)
		for val := range output {
			checks[int(val)] = true
		}

		for _, check := range checks {
			if !check {
				t.Fatal("not all data has been processed")
			}
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

		output := MapSync(func(val int) int32 {
			return int32(val)
		}, input)

		max := int32(-1)
		for val := range output {
			if val > max {
				max = val
			} else {
				t.Fatal("the data order is broken")
			}
		}

		if max < dataCount-1 {
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

		output := MapSequential(func(val int) int32 {
			return int32(val)
		}, input)

		max := int32(-1)
		for val := range output {
			if val > max {
				max = val
			} else {
				t.Fatal("the data order is broken")
			}
		}

		if max < dataCount-1 {
			t.Fatal("not all data has been processed")
		}
	})
}
