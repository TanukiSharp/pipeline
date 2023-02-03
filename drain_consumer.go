package pipeline

type DrainConsumer[T any] struct {
}

func NewDrainConsumer[T any]() *DrainConsumer[T] {
	return &DrainConsumer[T]{}
}

func (c *DrainConsumer[T]) Consume(input <-chan T) {
	if input == nil {
		panic("argument 'input' is mandatory")
	}

	for range input {
	}
}
