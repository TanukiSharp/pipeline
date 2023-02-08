package pipeline

type DelegateConsumer[T any] struct {
	name           string
	processingFunc func(input T) bool
}

var _ Consumer[any] = &DelegateConsumer[any]{}

func NewDelegateConsumer[T any](processingFunc func(input T) bool) *DelegateConsumer[T] {
	if processingFunc == nil {
		panic("argument 'processingFunc' is mandatory")
	}

	return &DelegateConsumer[T]{
		processingFunc: processingFunc,
	}
}

func (block *DelegateConsumer[T]) GetName() string {
	return block.name
}

func (block *DelegateConsumer[T]) SetName(name string) *DelegateConsumer[T] {
	block.name = name
	return block
}

func (block *DelegateConsumer[T]) Consume(input <-chan T) UnlinkFunc {
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
