// Package tests provides helpers to test the module functions.
package test

import (
	"fmt"
	"sync"
	"testing"

	"golang.org/x/exp/constraints"
)

// Consumer represent some kind of consumers.
func Consumer[T any](in <-chan T) {
	go func() {
		for {
			if _, ok := <-in; !ok {
				break
			}
		}
	}()
}

// Generator returns channel with data stream.
func Generator[T constraints.Integer](from T, to T, capacity int) <-chan T {
	if to < from {
		from, to = to, from
	}
	out := make(chan T, capacity)
	go func() {
		for i := T(from); i < to; i++ {
			out <- i
		}
		close(out)
	}()
	return out
}

// Wait lock current goroutine until the input won't be closed.
func Wait[T any](in <-chan T) <-chan struct{} {
	out := make(chan struct{})
	go func() {
		for {
			if _, ok := <-in; !ok {
				out <- struct{}{}
				break
			}
		}
	}()
	return out
}

// WaitError lock current goroutine until error is in the channel or the channel won't be closed.
func WaitError(ein <-chan error) <-chan error {
	out := make(chan error)
	go func() {
		for {
			if err, ok := <-ein; ok {
				out <- err
				break
			} else {
				out <- nil
				break
			}
		}
	}()
	return out
}

// Assert checks statement for every item.
func Assert[T constraints.Ordered](name string, in <-chan T, ein chan error, stmt func(data T) error) (<-chan T, chan error) {
	wg := sync.WaitGroup{}
	out := make(chan T, cap(in))
	eout := make(chan error, cap(ein))

	wg.Add(1)
	go func() {
		for {
			if data, ok := <-in; ok {
				if err := stmt(data); err != nil {
					eout <- fmt.Errorf("%s: %v", name, err)
				}
				out <- data
			} else {
				close(out)
				wg.Done()
				break
			}
		}
	}()

	go func() {
		for {
			if err, ok := <-ein; ok {
				eout <- err
			} else {
				wg.Wait()
				close(eout)
				break
			}
		}
	}()

	return out, eout
}

// AssertCount after the input closes checks if count of total elements passed through the input if equal the value.
func AssertCount[T constraints.Ordered](name string, in <-chan T, ein chan error, count int) (<-chan T, chan error) {
	wg := sync.WaitGroup{}
	out := make(chan T, cap(in))
	eout := make(chan error, cap(ein))

	items := 0
	wg.Add(1)
	go func() {
		for {
			if data, ok := <-in; ok {
				items++
				out <- data
			} else {
				close(out)
				wg.Done()
				break
			}
		}
	}()

	go func() {
		for {
			if err, ok := <-ein; ok {
				eout <- err
			} else {
				wg.Wait()
				if items != count {
					eout <- fmt.Errorf("%s: expected %d items, got %d", name, count, items)
				}
				close(eout)
				break
			}
		}
	}()

	return out, eout
}

// AssertOrderAsc checks if every item in ascending order.
func AssertOrderAsc[T constraints.Ordered](name string, in <-chan T, ein chan error) (<-chan T, chan error) {
	var max *T
	return Assert(name, in, ein, func(data T) error {
		if max == nil {
			max = &data
		} else {
			if data >= *max {
				max = &data
			} else {
				return fmt.Errorf("data order is broken: prev %v, current %v", *max, data)
			}
		}
		return nil
	})
}

// AssertBool checks if every item meets validation function.
func AssertBool[T constraints.Ordered](name string, in <-chan T, ein chan error, validator func(data T) bool) (<-chan T, chan error) {
	return Assert(name, in, ein, func(data T) error {
		if validator(data) {
			return nil
		}
		return fmt.Errorf("got %v", data)
	})
}

// Suit wraps to use asserts.
func Suit[T any](t *testing.T, tester func(epipe chan error) ([]<-chan T, chan error)) {
	const errorPipeCapacity = 1024
	epipe := make(chan error, errorPipeCapacity)
	pipes, errs := tester(epipe)

	go func() {
		wg := sync.WaitGroup{}
		for _, pipe := range pipes {
			pipe := pipe
			wg.Add(1)
			go func() {
				<-Wait(pipe)
				wg.Done()
			}()
		}
		wg.Wait()
		close(epipe)
	}()

	if err := <-WaitError(errs); err != nil {
		t.Fatal(err)
	}
}
