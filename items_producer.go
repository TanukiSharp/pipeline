package pipeline

import (
	"context"
	"fmt"
)

type ItemsProducer[T any] struct {
	name   string
	ctx    context.Context
	values []T
}

var _ Producer[any] = &ItemsProducer[any]{}

func NewItemsProducer[T any](ctx context.Context, values []T) *ItemsProducer[T] {
	fmt.Printf("[NewItemsProducer()] begin\n")
	defer fmt.Printf("[NewItemsProducer()] end\n")

	if ctx == nil {
		panic("argument 'ctx' is mandatory")
	}
	if values == nil {
		panic("argument 'values' is mandatory")
	}

	fmt.Printf("[NewItemsProducer()] values: %v\n", values)

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
	fmt.Printf("[ItemsProducer.Produce() %s] begin\n", block.GetName())
	defer fmt.Printf("[ItemsProducer.Produce() %s] end\n", block.GetName())

	output := make(chan T)
	fmt.Printf("[ItemsProducer.Produce() %s] output channel created\n", block.GetName())

	go func() {
		fmt.Printf("[ItemsProducer.Produce() %s] [goroutine] begin\n", block.GetName())
		defer fmt.Printf("[ItemsProducer.Produce() %s] [goroutine] end\n", block.GetName())

		defer func() {
			fmt.Printf("[ItemsProducer.Produce() %s] [goroutine] output channel closing\n", block.GetName())
			close(output)
			fmt.Printf("[ItemsProducer.Produce() %s] [goroutine] output channel closed\n", block.GetName())
		}()

		fmt.Printf("[ItemsProducer.Produce() %s] [goroutine] for loop begin (values: %v)\n", block.GetName(), block.values)
		for _, n := range block.values {
			fmt.Printf("[ItemsProducer.Produce() %s] [goroutine] for loop pull n = %v from values\n", block.GetName(), n)
			select {
			case output <- n:
				fmt.Printf("[ItemsProducer.Produce() %s] [goroutine] for loop, wrote n = %v to output channel\n", block.GetName(), n)
			case <-block.ctx.Done():
				fmt.Printf("[ItemsProducer.Produce() %s] [goroutine] for loop, context is done\n", block.GetName())
				return
			}
		}
		fmt.Printf("[ItemsProducer.Produce() %s] [goroutine] for loop end\n", block.GetName())
	}()

	return output
}

func (block *ItemsProducer[T]) LinkTo(consumer Consumer[T]) UnlinkFunc {
	fmt.Printf("[ItemsProducer.LinkTo() %s] begin\n", block.GetName())
	defer fmt.Printf("[ItemsProducer.LinkTo() %s] end\n", block.GetName())

	return consumer.Consume(block.Produce())
}
