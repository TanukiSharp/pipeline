package pipeline

import (
	"context"
	"sync"
)

NEED TO AVOID CLEAN FAN-IN AND FAN-OUT FUNCTIONS, AND USE THEM TO IMPLEMENT BLOCKS

type DelegatePropagatorBlock[TInput, TOutput any] struct {
	sync.Mutex
	name          string
	ctx           context.Context
	inputs        inputsSet[TInput, TOutput]
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

	unlink, isNew := block.inputs.set(input)

	if isNew == false {
		block.consumeCore()
	}

	return unlink
}

func (block *DelegatePropagatorBlock[TTInput, TOutput]) Produce() <-chan TOutput {
	return block.output
}

func (block *DelegatePropagatorBlock[TInput, TOutput]) LinkTo(target TargetBlock[TOutput]) UnlinkFunc {
	return target.Consume(block.Produce())
}
