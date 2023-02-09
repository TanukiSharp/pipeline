package pipeline

import "context"

type DelegateSourceBlock[T any] struct {
	name        string
	ctx         context.Context
	factoryFunc func() (T, bool)
}

var _ SourceBlock[any] = &DelegateSourceBlock[any]{}

func NewDelegateSourceBlock[T any](ctx context.Context, factoryFunc func() (T, bool)) *DelegateSourceBlock[T] {
	if ctx == nil {
		panic("argument 'ctx' is mandatory")
	}
	if factoryFunc == nil {
		panic("argument 'factoryFunc' is mandatory")
	}

	return &DelegateSourceBlock[T]{
		ctx:         ctx,
		factoryFunc: factoryFunc,
	}
}

func (block *DelegateSourceBlock[T]) GetName() string {
	return block.name
}

func (block *DelegateSourceBlock[T]) SetName(name string) *DelegateSourceBlock[T] {
	block.name = name
	return block
}

func (block *DelegateSourceBlock[T]) Produce() <-chan T {
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

func (block *DelegateSourceBlock[T]) LinkTo(target TargetBlock[T]) UnlinkFunc {
	return target.Consume(block.Produce())
}
