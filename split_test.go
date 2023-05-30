package pipe

import (
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
