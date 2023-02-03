package pipeline

type DelegateConsumer[T any] struct {
	processingFunc func(input T) bool
}

func NewDelegateConsumer[T any](processingFunc func(input T) bool) *DelegateConsumer[T] {
	if processingFunc == nil {
		panic("argument 'processingFunc' is mandatory")
	}

	return &DelegateConsumer[T]{
		processingFunc: processingFunc,
	}
}

func (c *DelegateConsumer[T]) Consume(input <-chan T) {
	if input == nil {
		panic("argument 'input' is mandatory")
	}

	for n := range input {
		if c.processingFunc(n) == false {
			break
		}
	}
}
