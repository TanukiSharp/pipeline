package pipeline

import "fmt"

type DrainConsumer[T any] struct {
	name string
}

var _ Consumer[any] = &DrainConsumer[any]{}

func NewDrainConsumer[T any]() *DrainConsumer[T] {
	fmt.Printf("[NewDrainConsumer()] begin\n")
	defer fmt.Printf("[NewDrainConsumer()] end\n")

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
	fmt.Printf("[DrainConsumer.Consume() %s] begin\n", block.GetName())
	defer fmt.Printf("[DrainConsumer.Consume() %s] end\n", block.GetName())

	if input == nil {
		panic("argument 'input' is mandatory")
	}

	for n := range input {
		fmt.Printf("[DrainConsumer.Consume() %s] for loop, read n = %v from input channel\n", block.GetName(), n)
	}

	return func() {}
}
