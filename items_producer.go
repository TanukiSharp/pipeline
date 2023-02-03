package pipeline

type ItemsProducer[T any] struct {
	done   <-chan struct{}
	values []T
}

func NewItemsProducer[T any](pipeline *Pipeline[T], values []T) *ItemsProducer[T] {
	if pipeline == nil {
		panic("argument 'pipeline' is mandatory")
	}
	if values == nil {
		panic("argument 'values' is mandatory")
	}

	return &ItemsProducer[T]{
		done:   pipeline.GetDone(),
		values: values,
	}
}

func (g *ItemsProducer[T]) Produce() <-chan T {
	output := make(chan T)

	go func() {
		defer close(output)
		for _, n := range g.values {
			select {
			case output <- n:
			case <-g.done:
				return
			}
		}
	}()

	return output
}
