// - Sequential functions guarantee that only one item from input channel will be processed at a time, and potentially thread-safe to work with objects outside from handler function.
// - Not sequential functions processes all items at a same time and have no guarantee that output order will be follow input order.
// - Every function produces output channel with the same capacity as input channel.
// - Every function will close output channels if all input channels are closed.
package pipe

import "sync"

// Map transforms input chanel data and writes it into the output channel.
// If input channel has no buffer it acts like MapSequential.
// If input channel has buffer it has no guarantee that output order will follow input order.
// Output channel will have same capacity as input channel.
//
//   - Split - Потребляет 1 значение и дублирует его в каждый канал. Если входной закрывается - закрывает все выходящие.
//   - Route - Потребляет 1 значение и определяет в какой канал его направить. Если входной закрывается - закрывает все выходящие.
//   - Join - Потребляет несколько каналов и пишет в один. Если все входные закрываются - закрывает выходящий.
//   - Merge - Потребляет по одному значению из каждого входного канала и пишет в один. Если один из каналов закрывается - закрывает выходящий.
//   - Generate - Потребляет 1 значение и возвращает несколько в один канал. Если входной закрывается - закрывает выходящий.
//
// Неудачные концепции:
//
//   - Validate - Аналогичен Route, нет необходимости
//   - Repeater - Потребляет 1 значение, выдаёт его на выход и решает, повторить ли его на входной канал (нарушает паттерн потока данных, аналогичен генератору по своей сути)
func Map[Tin, Tout any](in <-chan Tin, mapper func(Tin) Tout) <-chan Tout {
	out := make(chan Tout, cap(in))
	wg := sync.WaitGroup{}
	go func() {
		for {
			if in, ok := <-in; ok {
				wg.Add(1)
				go func() {
					out <- mapper(in)
					wg.Done()
				}()
			} else {
				wg.Wait()
				close(out)
				break
			}
		}
	}()
	return out
}

func MapSequential[Tin, Tout any](in <-chan Tin, mapper func(Tin) Tout) <-chan Tout {
	out := make(chan Tout, cap(in))
	go func() {
		for {
			if in, ok := <-in; ok {
				out <- mapper(in)
			} else {
				close(out)
				break
			}
		}
	}()
	return out
}

func Filter[T any](in <-chan T, filter func(T) bool) <-chan T {
	out := make(chan T, cap(in))
	wg := sync.WaitGroup{}
	go func() {
		for {
			if in, ok := <-in; ok {
				wg.Add(1)
				go func() {
					if filter(in) {
						out <- in
					}
					wg.Done()
				}()
			} else {
				wg.Wait()
				close(out)
				break
			}
		}
	}()
	return out
}

func FilterSequential[T any](in <-chan T, filter func(T) bool) <-chan T {
	out := make(chan T, cap(in))
	go func() {
		for {
			if in, ok := <-in; ok {
				if filter(in) {
					out <- in
				}
			} else {
				close(out)
				break
			}
		}
	}()
	return out
}

// Watch calls callback for every item in sequence without ordering.
func Watch[T any](in <-chan T, callback func(T)) <-chan T {
	out := make(chan T, cap(in))
	wg := sync.WaitGroup{}
	go func() {
		for {
			if in, ok := <-in; ok {
				wg.Add(1)
				go func() {
					callback(in)
					out <- in
					wg.Done()
				}()
			} else {
				wg.Wait()
				close(out)
				break
			}
		}
	}()
	return out
}

// WatchSequential calls callback for every item in sequence with input order.
func WatchSequential[T any](in <-chan T, callback func(T)) <-chan T {
	out := make(chan T, cap(in))
	go func() {
		for {
			if in, ok := <-in; ok {
				callback(in)
				out <- in
			} else {
				close(out)
				break
			}
		}
	}()
	return out
}

// Wait is terminate function. It will block current goroutine until input channel is closed.
func Wait[T any](in <-chan T) {
	for {
		if _, ok := <-in; !ok {
			break
		}
	}
}

// func Generate[T any](generator func() (T, bool)) <-chan T {
// }
