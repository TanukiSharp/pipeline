package pipeline

type DrainTargetBlock[T any] struct {
	name string
}

var _ TargetBlock[any] = &DrainTargetBlock[any]{}

func NewDrainTargetBlock[T any]() *DrainTargetBlock[T] {
	return &DrainTargetBlock[T]{}
}

func (block *DrainTargetBlock[T]) GetName() string {
	return block.name
}

func (block *DrainTargetBlock[T]) SetName(name string) *DrainTargetBlock[T] {
	block.name = name
	return block
}

func (block *DrainTargetBlock[T]) Consume(input <-chan T) UnlinkFunc {
	if input == nil {
		panic("argument 'input' is mandatory")
	}

	for range input {
	}

	return func() {}
}
