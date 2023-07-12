package pipeline

import (
	"context"
)

type ItemsSourceBlock[T any] struct {
	name   string
	ctx    context.Context
	values []T
	output chan T
}

var _ SourceBlock[any] = &ItemsSourceBlock[any]{}

func NewItemsSourceBlock[T any](ctx context.Context, values []T) *ItemsSourceBlock[T] {
	if ctx == nil {
		panic("argument 'ctx' is mandatory")
	}
	if values == nil {
		panic("argument 'values' is mandatory")
	}

	return &ItemsSourceBlock[T]{
		ctx:    ctx,
		values: values,
	}
}

func (block *ItemsSourceBlock[T]) GetName() string {
	return block.name
}

func (block *ItemsSourceBlock[T]) SetName(name string) *ItemsSourceBlock[T] {
	block.name = name
	return block
}

func (block *ItemsSourceBlock[T]) Produce() <-chan T {
	if block.output != nil {
		return block.output
	}

	block.output = make(chan T)

	go func() {
		defer close(block.output)

		for _, n := range block.values {
			select {
			case block.output <- n:
			case <-block.ctx.Done():
				return
			}
		}
	}()

	return block.output
}

func (block *ItemsSourceBlock[T]) LinkTo(target TargetBlock[T]) UnlinkFunc {
	return target.Consume(block.Produce())
}
