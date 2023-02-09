package pipeline

import (
	"context"
)

type DelegateSubject[TInput, TOutput any] struct {
	name          string
	ctx           context.Context
	output        chan TOutput
	transformFunc func(input TInput) TOutput
}

var _ Subject[any, any] = &DelegateSubject[any, any]{}

func NewDelegateSubject[TInput, TOutput any](ctx context.Context, transformFunc func(input TInput) TOutput) *DelegateSubject[TInput, TOutput] {
	if ctx == nil {
		panic("argument 'ctx' is mandatory")
	}
	if transformFunc == nil {
		panic("argument 'transformFunc' is mandatory")
	}

	return &DelegateSubject[TInput, TOutput]{
		ctx:           ctx,
		transformFunc: transformFunc,
	}
}

func (block *DelegateSubject[TInput, TOutput]) GetName() string {
	return block.name
}

func (block *DelegateSubject[TInput, TOutput]) SetName(name string) *DelegateSubject[TInput, TOutput] {
	block.name = name
	return block
}

func (block *DelegateSubject[TInput, TOutput]) Consume(input <-chan TInput) UnlinkFunc {
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

func (block *DelegateSubject[TTInput, TOutput]) Produce() <-chan TOutput {
	return block.output
}

func (block *DelegateSubject[TInput, TOutput]) LinkTo(consumer Consumer[TOutput]) UnlinkFunc {
	return consumer.Consume(block.Produce())
}
