package pipeline

type DelegateProducer[T any] struct {
	done        <-chan struct{}
	factoryFunc func() (T, bool)
}

func NewDelegateProducer[T any](pipeline *Pipeline[T], factoryFunc func() (T, bool)) *DelegateProducer[T] {
	if pipeline == nil {
		panic("argument 'pipeline' is mandatory")
	}
	if factoryFunc == nil {
		panic("argument 'factoryFunc' is mandatory")
	}

	return &DelegateProducer[T]{
		done:        pipeline.GetDone(),
		factoryFunc: factoryFunc,
	}
}

func (p *DelegateProducer[T]) Produce() <-chan T {
	output := make(chan T)

	go func() {
		defer close(output)
		for {
			item, hasItem := p.factoryFunc()
			if hasItem == false {
				break
			}
			select {
			case output <- item:
			case <-p.done:
				return
			}
		}
	}()

	return output
}
