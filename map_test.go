package pipe

import (
	"testing"

	"github.com/msacore/pipe/test"
)

func TestMap(t *testing.T) {
	t.Run("Parallel", func(t *testing.T) {
		test.Suit(t, func(epipe chan error) ([]<-chan float32, chan error) {
			pipe := test.Generator(0, 64, 16)

			pipe2 := Map(func(val int) float32 {
				return float32(val)
			}, pipe)

			pipe2, epipe = test.AssertCount("count", pipe2, epipe, 64)

			return []<-chan float32{pipe2}, epipe
		})
	})

	t.Run("Sync", func(t *testing.T) {
		test.Suit(t, func(epipe chan error) ([]<-chan float32, chan error) {
			pipe := test.Generator(0, 64, 16)

			pipe2 := MapSync(func(val int) float32 {
				return float32(val)
			}, pipe)

			pipe2, epipe = test.AssertOrderAsc("ordering", pipe2, epipe)
			pipe2, epipe = test.AssertCount("count", pipe2, epipe, 64)

			return []<-chan float32{pipe2}, epipe
		})
	})

	t.Run("Sequential", func(t *testing.T) {
		test.Suit(t, func(epipe chan error) ([]<-chan float32, chan error) {
			pipe := test.Generator(0, 64, 16)

			pipe2 := MapSequential(func(val int) float32 {
				return float32(val)
			}, pipe)

			pipe2, epipe = test.AssertOrderAsc("ordering", pipe2, epipe)
			pipe2, epipe = test.AssertCount("count", pipe2, epipe, 64)

			return []<-chan float32{pipe2}, epipe
		})
	})
}
