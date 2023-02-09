package pipeline

import (
	"context"
	"sync"
)

type MergePropagatorBlock[TInput, TOutput any] struct {
	name        string
	ctx         context.Context
	propagators []PropagatorBlock[TInput, TOutput]
}

var _ PropagatorBlock[any, any] = &MergePropagatorBlock[any, any]{}

func NewMergePropagatorBlockWithFactory[TInput, TOutput any](ctx context.Context, howMany int, factoryFunc func(i int) PropagatorBlock[TInput, TOutput]) *MergePropagatorBlock[TInput, TOutput] {
	if ctx == nil {
		panic("argument 'ctx' is mandatory")
	}
	if howMany <= 0 {
		panic("argument 'howMany' must be positive")
	}
	if factoryFunc == nil {
		panic("argument 'factoryFunc' is mandatory")
	}

	propagators := make([]PropagatorBlock[TInput, TOutput], howMany)

	for i := 0; i < howMany; i++ {
		propagators[i] = factoryFunc(i)
	}

	return NewMergePropagatorBlockWithInstances(ctx, propagators...)
}

func NewMergePropagatorBlockWithInstances[TInput, TOutput any](ctx context.Context, propagators ...PropagatorBlock[TInput, TOutput]) *MergePropagatorBlock[TInput, TOutput] {
	if ctx == nil {
		panic("argument 'ctx' is mandatory")
	}
	if propagators == nil {
		panic("argument 'propagators' is mandatory")
	}
	if len(propagators) == 0 {
		panic("argument 'propagators' must contain at least one element")
	}

	return &MergePropagatorBlock[TInput, TOutput]{
		ctx:         ctx,
		propagators: propagators,
	}
}

func (block *MergePropagatorBlock[TInput, TOutput]) GetName() string {
	return block.name
}

func (block *MergePropagatorBlock[TInput, TOutput]) SetName(name string) *MergePropagatorBlock[TInput, TOutput] {
	block.name = name
	return block
}

func (block *MergePropagatorBlock[TInput, TOutput]) Consume(input <-chan TInput) UnlinkFunc {
	if input == nil {
		panic("argument 'input' is mandatory")
	}

	unlinkFuncs := []func(){}

	for _, p := range block.propagators {
		unlinkFunc := p.Consume(input)
		unlinkFuncs = append(unlinkFuncs, unlinkFunc)
	}

	return func() {
		for _, unlinkFunc := range unlinkFuncs {
			unlinkFunc()
		}
	}
}

func (block *MergePropagatorBlock[TInput, TOutput]) Produce() <-chan TOutput {
	var wg sync.WaitGroup
	output := make(chan TOutput)

	outputFunc := func(c <-chan TOutput) {
		defer wg.Done()
		for n := range c {
			select {
			case output <- n:
			case <-block.ctx.Done():
				return
			}
		}
	}

	wg.Add(len(block.propagators))

	for _, p := range block.propagators {
		go outputFunc(p.Produce())
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

func (block *MergePropagatorBlock[TInput, TOutput]) LinkTo(target TargetBlock[TOutput]) UnlinkFunc {
	return target.Consume(block.Produce())
}
