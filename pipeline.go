package pipeline

import "fmt"

type Pipeline[T any] struct {
	done     chan struct{}
	source   Producer[T]
	subjects []Subject[T]
	sink     Consumer[T]
}

func NewPipeline[T any]() *Pipeline[T] {
	return &Pipeline[T]{
		done:     make(chan struct{}),
		subjects: []Subject[T]{},
	}
}

func (p *Pipeline[T]) RegisterSource(source Producer[T]) *Pipeline[T] {
	p.source = source
	return p
}

func (p *Pipeline[T]) RegisterSink(sink Consumer[T]) *Pipeline[T] {
	p.sink = sink
	return p
}

func (p *Pipeline[T]) AddSubject(subject Subject[T]) *Pipeline[T] {
	p.subjects = append(p.subjects, subject)
	return p
}

func (p *Pipeline[T]) Run() error {
	if p.source == nil {
		return fmt.Errorf("a source must be registered")
	}

	if p.sink == nil {
		return fmt.Errorf("a sink must be registered")
	}

	c := p.source.Produce()

	for _, s := range p.subjects {
		s.Consume(c)
		c = s.Produce()
	}

	p.sink.Consume(c)

	return nil
}

func (p *Pipeline[T]) GetDone() <-chan struct{} {
	return p.done
}

func (p *Pipeline[T]) Cancel() {
	close(p.done)
}
