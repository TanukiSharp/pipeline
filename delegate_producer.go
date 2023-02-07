package pipeline

import "context"

type DelegateProducer[T any] struct {
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

func (p *DelegateProducer[T]) Produce() <-chan T {
	output := make(chan T)

	go func() {
		defer close(output)
		for {
			item, hasItem := p.factoryFunc()
			if hasItem == false {
				break
			}
			select {
			case output <- item:
			case <-p.ctx.Done():
				return
			}
		}
	}()

	return output
}

func (p *DelegateProducer[T]) LinkTo(consumer Consumer[T]) func() {
	return consumer.Consume(p.Produce())
}
