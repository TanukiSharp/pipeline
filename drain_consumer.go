package pipeline

type DrainConsumer[T any] struct {
	name string
}

var _ Consumer[any] = &DrainConsumer[any]{}

func NewDrainConsumer[T any]() *DrainConsumer[T] {
	return &DrainConsumer[T]{}
}

func (block *DrainConsumer[T]) GetName() string {
	return block.name
}

func (block *DrainConsumer[T]) SetName(name string) *DrainConsumer[T] {
	block.name = name
	return block
}

func (block *DrainConsumer[T]) Consume(input <-chan T) UnlinkFunc {
	if input == nil {
		panic("argument 'input' is mandatory")
	}

	for range input {
	}

	return func() {}
}
