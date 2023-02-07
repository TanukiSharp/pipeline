package pipeline

import "context"

type ItemsProducer[T any] struct {
	ctx    context.Context
	values []T
}

var _ Producer[any] = &ItemsProducer[any]{}

func NewItemsProducer[T any](ctx context.Context, values []T) *ItemsProducer[T] {
	if ctx == nil {
		panic("argument 'ctx' is mandatory")
	}
	if values == nil {
		panic("argument 'values' is mandatory")
	}

	return &ItemsProducer[T]{
		ctx:    ctx,
		values: values,
	}
}

func (g *ItemsProducer[T]) Produce() <-chan T {
	output := make(chan T)

	go func() {
		defer close(output)
		for _, n := range g.values {
			select {
			case output <- n:
			case <-g.ctx.Done():
				return
			}
		}
	}()

	return output
}

func (g *ItemsProducer[T]) LinkTo(consumer Consumer[T]) func() {
	return consumer.Consume(g.Produce())
}
