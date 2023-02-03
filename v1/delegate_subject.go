package pipeline

type DelegateSubject[T any] struct {
	done          <-chan struct{}
	output        chan T
	transformFunc func(input T) T
}

func NewDelegateSubject[T any](pipeline *Pipeline[T], transformFunc func(input T) T) *DelegateSubject[T] {
	if pipeline == nil {
		panic("argument 'pipeline' is mandatory")
	}
	if transformFunc == nil {
		panic("argument 'transformFunc' is mandatory")
	}

	return &DelegateSubject[T]{
		done:          pipeline.GetDone(),
		transformFunc: transformFunc,
	}
}

func (s *DelegateSubject[T]) Consume(input <-chan T) {
	if input == nil {
		panic("argument 'input' is mandatory")
	}

	s.output = make(chan T)

	go func() {
		defer close(s.output)
		for n := range input {
			select {
			case s.output <- s.transformFunc(n):
			case <-s.done:
				return
			}
		}
	}()
}

func (s *DelegateSubject[T]) Produce() <-chan T {
	return s.output
}
