package pipeline

import (
	"context"
)

type ItemsProducer[T any] struct {
	name   string
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

func (block *ItemsProducer[T]) GetName() string {
	return block.name
}

func (block *ItemsProducer[T]) SetName(name string) *ItemsProducer[T] {
	block.name = name
	return block
}

func (block *ItemsProducer[T]) Produce() <-chan T {
	output := make(chan T)

	go func() {
		defer close(output)

		for _, n := range block.values {
			select {
			case output <- n:
			case <-block.ctx.Done():
				return
			}
		}
	}()

	return output
}

func (block *ItemsProducer[T]) LinkTo(consumer Consumer[T]) UnlinkFunc {
	return consumer.Consume(block.Produce())
}
