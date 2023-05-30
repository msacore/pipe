package pipe

import (
	"fmt"
	"testing"

	"github.com/msacore/pipe/test"
)

func TestSplit(t *testing.T) {
	const channelsCount = 10

	t.Run("Parallel", func(t *testing.T) {
		test.Suit(t, func(epipe chan error) ([]<-chan int, chan error) {
			pipe := test.Generator(0, 64, 16)

			pipes := Split(channelsCount, pipe)
			for i := 0; i < channelsCount; i++ {
				pipes[i], epipe = test.AssertCount(fmt.Sprintf("count %d", i), pipes[i], epipe, 64)
			}

			return pipes, epipe
		})
	})

	t.Run("Sync", func(t *testing.T) {
		test.Suit(t, func(epipe chan error) ([]<-chan int, chan error) {
			pipe := test.Generator(0, 64, 16)

			pipes := SplitSync(channelsCount, pipe)
			for i := 0; i < channelsCount; i++ {
				pipes[i], epipe = test.AssertCount(fmt.Sprintf("count %d", i), pipes[i], epipe, 64)
				pipes[i], epipe = test.AssertOrderAsc(fmt.Sprintf("ordering %d", i), pipes[i], epipe)
			}

			return pipes, epipe
		})
	})

	t.Run("Sequential", func(t *testing.T) {
		test.Suit(t, func(epipe chan error) ([]<-chan int, chan error) {
			pipe := test.Generator(0, 64, 16)

			pipes := SplitSequential(channelsCount, pipe)
			for i := 0; i < channelsCount; i++ {
				pipes[i], epipe = test.AssertCount(fmt.Sprintf("count %d", i), pipes[i], epipe, 64)
				pipes[i], epipe = test.AssertOrderAsc(fmt.Sprintf("ordering %d", i), pipes[i], epipe)
			}

			return pipes, epipe
		})
	})
}
