package pipe

import (
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	t.Run("Wait", func(t *testing.T) {
		input := make(chan struct{}, 1)

		go func() {
			<-time.After(time.Millisecond * 5)
			input <- struct{}{}
			close(input)
		}()

		start := time.Now()
		<-Wait(input)

		if time.Since(start).Milliseconds() < 4 {
			t.Fatal("too fast")
		}
		if time.Since(start).Milliseconds() > 6 {
			t.Fatal("too slow")
		}
	})

	t.Run("WaitAll", func(t *testing.T) {
		input1 := make(chan struct{}, 1)
		input2 := make(chan struct{}, 1)

		go func() {
			<-time.After(time.Millisecond * 5)
			input1 <- struct{}{}
			close(input1)
			<-time.After(time.Millisecond * 3)
			input2 <- struct{}{}
			close(input2)
		}()

		start := time.Now()
		<-WaitAll(input1, input2)

		if time.Since(start).Milliseconds() < 7 {
			t.Fatal("too fast")
		}
		if time.Since(start).Milliseconds() > 9 {
			t.Fatal("too slow")
		}
	})

	t.Run("WaitAny", func(t *testing.T) {
		input1 := make(chan struct{}, 1)
		input2 := make(chan struct{}, 1)

		go func() {
			<-time.After(time.Millisecond * 5)
			input1 <- struct{}{}
			close(input1)
			<-time.After(time.Millisecond * 3)
			input2 <- struct{}{}
			close(input2)
		}()

		start := time.Now()
		<-WaitAny(input1, input2)

		if time.Since(start).Milliseconds() < 4 {
			t.Fatal("too fast")
		}
		if time.Since(start).Milliseconds() > 6 {
			t.Fatal("too slow")
		}
	})
}
