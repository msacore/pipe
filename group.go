package pipe

type Group[T any] []chan T

func (g Group[T]) Readers() GroupReaders[T] {
	channels := make([]<-chan T, len(g))
	for i := range g {
		channels[i] = g[i]
	}
	return channels
}

func (g Group[T]) Writers() GroupWriters[T] {
	channels := make([]chan<- T, len(g))
	for i := range g {
		channels[i] = g[i]
	}
	return channels
}

func (g Group[T]) As1() (ch1 chan T) {
	if len(g) != 1 {
		panic("group length must be 1")
	}
	return g[0]
}

func (g Group[T]) As2() (ch1, ch2 chan T) {
	if len(g) != 2 {
		panic("group length must be 2")
	}
	return g[0], g[1]
}

func (g Group[T]) As3() (ch1, ch2, ch3 chan T) {
	if len(g) != 3 {
		panic("group length must be 3")
	}
	return g[0], g[1], g[2]
}

type GroupReaders[T any] []<-chan T

type GroupWriters[T any] []chan<- T
