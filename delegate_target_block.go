package pipeline

type DelegateTargetBlock[T any] struct {
	name           string
	processingFunc func(input T) bool
}

var _ TargetBlock[any] = &DelegateTargetBlock[any]{}

func NewDelegateTargetBlock[T any](processingFunc func(input T) bool) *DelegateTargetBlock[T] {
	if processingFunc == nil {
		panic("argument 'processingFunc' is mandatory")
	}

	return &DelegateTargetBlock[T]{
		processingFunc: processingFunc,
	}
}

func (block *DelegateTargetBlock[T]) GetName() string {
	return block.name
}

func (block *DelegateTargetBlock[T]) SetName(name string) *DelegateTargetBlock[T] {
	block.name = name
	return block
}

func (block *DelegateTargetBlock[T]) Consume(input <-chan T) UnlinkFunc {
	if input == nil {
		panic("argument 'input' is mandatory")
	}

	for n := range input {
		if block.processingFunc(n) == false {
			break
		}
	}

	return func() {}
}
