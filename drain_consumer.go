package pipeline

type DrainConsumer[T any] struct {
}

var _ Consumer[any] = &DrainConsumer[any]{}

func NewDrainConsumer[T any]() *DrainConsumer[T] {
	return &DrainConsumer[T]{}
}

func (c *DrainConsumer[T]) Consume(input <-chan T) func() {
	if input == nil {
		panic("argument 'input' is mandatory")
	}

	for range input {
	}

	return func() {}
}
