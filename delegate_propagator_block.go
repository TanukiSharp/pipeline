package pipeline

import (
	"context"
)

type DelegatePropagatorBlock[TInput, TOutput any] struct {
	name          string
	ctx           context.Context
	output        chan TOutput
	transformFunc func(input TInput) TOutput
}

var _ PropagatorBlock[any, any] = &DelegatePropagatorBlock[any, any]{}

func NewDelegatePropagatorBlock[TInput, TOutput any](ctx context.Context, transformFunc func(input TInput) TOutput) *DelegatePropagatorBlock[TInput, TOutput] {
	if ctx == nil {
		panic("argument 'ctx' is mandatory")
	}
	if transformFunc == nil {
		panic("argument 'transformFunc' is mandatory")
	}

	return &DelegatePropagatorBlock[TInput, TOutput]{
		ctx:           ctx,
		transformFunc: transformFunc,
	}
}

func (block *DelegatePropagatorBlock[TInput, TOutput]) GetName() string {
	return block.name
}

func (block *DelegatePropagatorBlock[TInput, TOutput]) SetName(name string) *DelegatePropagatorBlock[TInput, TOutput] {
	block.name = name
	return block
}

func (block *DelegatePropagatorBlock[TInput, TOutput]) Consume(input <-chan TInput) UnlinkFunc {
	if input == nil {
		panic("argument 'input' is mandatory")
	}

	linkContext, unlink := context.WithCancel(context.Background())

	block.output = make(chan TOutput)

	go func() {
		defer close(block.output)

		for {
			select {
			case n, ok := <-input:
				if ok == false {
					return
				}

				select {
				case block.output <- block.transformFunc(n):
				case <-block.ctx.Done():
					return
				case <-linkContext.Done():
					return
				}
			case <-block.ctx.Done():
				return
			case <-linkContext.Done():
				return
			}
		}
	}()

	return func() {
		unlink()
	}
}

func (block *DelegatePropagatorBlock[TTInput, TOutput]) Produce() <-chan TOutput {
	return block.output
}

func (block *DelegatePropagatorBlock[TInput, TOutput]) LinkTo(target TargetBlock[TOutput]) UnlinkFunc {
	return target.Consume(block.Produce())
}
