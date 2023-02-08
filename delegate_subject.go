package pipeline

import (
	"context"
	"fmt"
)

type DelegateSubject[TInput, TOutput any] struct {
	name          string
	ctx           context.Context
	output        chan TOutput
	transformFunc func(input TInput) TOutput
}

var _ Subject[any, any] = &DelegateSubject[any, any]{}

func NewDelegateSubject[TInput any, TOutput any](ctx context.Context, transformFunc func(input TInput) TOutput) *DelegateSubject[TInput, TOutput] {
	fmt.Printf("[NewDelegateSubject()] begin\n")
	defer fmt.Printf("[NewDelegateSubject()] end\n")

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
	fmt.Printf("[DelegateSubject.Consume() %s] begin\n", block.GetName())
	defer fmt.Printf("[DelegateSubject.Consume() %s] end\n", block.GetName())

	if input == nil {
		panic("argument 'input' is mandatory")
	}

	isRunning, cancel := context.WithCancel(context.Background())

	block.output = make(chan TOutput)
	fmt.Printf("[DelegateSubject.Consume() %s] output channel created\n", block.GetName())

	go func() {
		fmt.Printf("[DelegateSubject.Consume() %s] [goroutine] begin\n", block.GetName())
		defer fmt.Printf("[DelegateSubject.Consume() %s] [goroutine] end\n", block.GetName())

		defer func() {
			fmt.Printf("[DelegateSubject.Consume() %s] [goroutine] output channel closing\n", block.GetName())
			close(block.output)
			fmt.Printf("[DelegateSubject.Consume() %s] [goroutine] output channel closed\n", block.GetName())
		}()

		fmt.Printf("[DelegateSubject.Consume() %s] [goroutine] for loop begin (input channel len: %d)\n", block.GetName(), len(input))

		defer fmt.Printf("[DelegateSubject.Consume() %s] [goroutine] for loop end\n", block.GetName())
		for {
			select {
			case n, ok := <-input:
				if ok == false {
					fmt.Printf("[DelegateSubject.Consume() %s] [goroutine] for loop, input channel closed\n", block.GetName())
					return // This breaks the 'case', not the 'for' loop.
				}
				fmt.Printf("[DelegateSubject.Consume() %s] [goroutine] for loop, pull n = %v from input channel\n", block.GetName(), n)
				select {
				case block.output <- block.transformFunc(n):
					fmt.Printf("[DelegateSubject.Consume() %s] [goroutine] for loop, wrote n = %v to output channel\n", block.GetName(), n)
				case <-block.ctx.Done():
					fmt.Printf("[DelegateSubject.Consume() %s] [goroutine] for loop, context is done\n", block.GetName())
				case <-isRunning.Done():
					fmt.Printf("[DelegateSubject.Consume() %s] [goroutine] for loop, link as been unlinked\n", block.GetName())
					return
				}
			case <-block.ctx.Done():
				fmt.Printf("[DelegateSubject.Consume() %s] [goroutine] for loop, context is done\n", block.GetName())
				return
			case <-isRunning.Done():
				fmt.Printf("[DelegateSubject.Consume() %s] [goroutine] for loop, link as been unlinked\n", block.GetName())
				return
			}
		}
	}()

	return func() {
		cancel()
	}
}

func (block *DelegateSubject[TTInput, TOutput]) Produce() <-chan TOutput {
	fmt.Printf("[DelegateSubject.Produce() %s] begin\n", block.GetName())
	defer fmt.Printf("[DelegateSubject.Produce() %s] end\n", block.GetName())

	return block.output
}

func (block *DelegateSubject[TInput, TOutput]) LinkTo(consumer Consumer[TOutput]) UnlinkFunc {
	fmt.Printf("[DelegateSubject.LinkTo() %s] begin\n", block.GetName())
	defer fmt.Printf("[DelegateSubject.LinkTo() %s] end\n", block.GetName())

	return consumer.Consume(block.Produce())
}
