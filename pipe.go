package pipe

// ReadWriter describes generic Pipe interface.
type ReadWriter[T any] interface {
	Writer[T]
	Reader[T]
}

// Writer describes restricted Pipe interface with only Write ability.
type Writer[T any] interface {
	Write(T)
}

// Reader describes restricted Pipe interface with only Read ability.
type Reader[T any] interface {
	Read() T
}

// Pipe implements generic Pipe interface.
type Pipe[T any] struct {
	ch chan T
}

// PipeWriter implements generic restricted Pipe interface with only Write ability.
type PipeWriter[T any] struct {
	pipe *Pipe[T]
}

// PipeReader implements generic restricted Pipe interface with only Read ability.
type PipeReader[T any] struct {
	pipe *Pipe[T]
}

// New makes a new Pipe with specific type.
// Pipes is a channels' wrappers without buffering.
// It's allows to write and read messages.
func New[T any]() ReadWriter[T] {
	return &Pipe[T]{
		ch: make(chan T),
	}
}

func (p *Pipe[T]) Write(msg T) {
	p.ch <- msg
}

func (p *Pipe[T]) Read() T {
	return <-p.ch
}

func (p *Pipe[T]) Writer() Writer[T] {
	return &PipeWriter[T]{
		pipe: p,
	}
}

func (p *Pipe[T]) Reader() Reader[T] {
	return &PipeReader[T]{
		pipe: p,
	}
}

func (p *PipeWriter[T]) Write(msg T) {
	p.pipe.Write(msg)
}

func (p *PipeReader[T]) Read() T {
	return p.pipe.Read()
}
