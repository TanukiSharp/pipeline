package pipeline

import "context"

type DelegateProducer[T any] struct {
	name        string
	ctx         context.Context
	factoryFunc func() (T, bool)
}

var _ Producer[any] = &DelegateProducer[any]{}

func NewDelegateProducer[T any](ctx context.Context, factoryFunc func() (T, bool)) *DelegateProducer[T] {
	if ctx == nil {
		panic("argument 'ctx' is mandatory")
	}
	if factoryFunc == nil {
		panic("argument 'factoryFunc' is mandatory")
	}

	return &DelegateProducer[T]{
		ctx:         ctx,
		factoryFunc: factoryFunc,
	}
}

func (block *DelegateProducer[T]) GetName() string {
	return block.name
}

func (block *DelegateProducer[T]) SetName(name string) *DelegateProducer[T] {
	block.name = name
	return block
}

func (block *DelegateProducer[T]) Produce() <-chan T {
	output := make(chan T)

	go func() {
		defer close(output)
		for {
			item, hasItem := block.factoryFunc()
			if hasItem == false {
				break
			}
			select {
			case output <- item:
			case <-block.ctx.Done():
				return
			}
		}
	}()

	return output
}

func (block *DelegateProducer[T]) LinkTo(consumer Consumer[T]) UnlinkFunc {
	return consumer.Consume(block.Produce())
}
