package pipe

import (
	"testing"

	"github.com/msacore/pipe/test"
)

func TestFilter(t *testing.T) {
	t.Run("Parallel", func(t *testing.T) {
		test.Suit(t, func(epipe chan error) ([]<-chan int, chan error) {
			pipe := test.Generator(0, 64, 16)

			pipe = Filter(func(val int) bool {
				return val%2 == 0
			}, pipe)

			pipe, epipe = test.AssertBool("validation", pipe, epipe, func(data int) bool { return data%2 == 0 })
			pipe, epipe = test.AssertCount("count", pipe, epipe, 32)

			return []<-chan int{pipe}, epipe
		})
	})

	t.Run("Sync", func(t *testing.T) {
		test.Suit(t, func(epipe chan error) ([]<-chan int, chan error) {
			pipe := test.Generator(0, 64, 16)

			pipe = FilterSync(func(val int) bool {
				return val%2 == 0
			}, pipe)

			pipe, epipe = test.AssertBool("validation", pipe, epipe, func(data int) bool { return data%2 == 0 })
			pipe, epipe = test.AssertOrderAsc("ordering", pipe, epipe)
			pipe, epipe = test.AssertCount("count", pipe, epipe, 32)

			return []<-chan int{pipe}, epipe
		})
	})

	t.Run("Sequential", func(t *testing.T) {
		test.Suit(t, func(epipe chan error) ([]<-chan int, chan error) {
			pipe := test.Generator(0, 64, 16)

			pipe = FilterSequential(func(val int) bool {
				return val%2 == 0
			}, pipe)

			pipe, epipe = test.AssertBool("validation", pipe, epipe, func(data int) bool { return data%2 == 0 })
			pipe, epipe = test.AssertOrderAsc("ordering", pipe, epipe)
			pipe, epipe = test.AssertCount("count", pipe, epipe, 32)

			return []<-chan int{pipe}, epipe
		})
	})
}
