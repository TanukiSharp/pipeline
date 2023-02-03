package pipeline

import "sync"

type MergeSubject[T any] struct {
	done     <-chan struct{}
	subjects []Subject[T]
}

func NewMergeSubjectWithFactory[T any](pipeline *Pipeline[T], howMany int, factoryFunc func() Subject[T]) *MergeSubject[T] {
	if pipeline == nil {
		panic("argument 'pipeline' is mandatory")
	}
	if factoryFunc == nil {
		panic("argument 'factoryFunc' is mandatory")
	}

	subjects := make([]Subject[T], howMany)

	for i := 0; i < howMany; i++ {
		subjects[i] = factoryFunc()
	}

	return NewMergeSubjectWithInstances(pipeline, subjects...)
}

func NewMergeSubjectWithInstances[T any](pipeline *Pipeline[T], subjects ...Subject[T]) *MergeSubject[T] {
	if pipeline == nil {
		panic("argument 'pipeline' is mandatory")
	}
	if subjects == nil {
		panic("argument 'subject' is mandatory")
	}
	if len(subjects) == 0 {
		panic("argument 'subjects' must contain at least one element")
	}

	return &MergeSubject[T]{
		done:     pipeline.GetDone(),
		subjects: subjects,
	}
}

func (c *MergeSubject[T]) Consume(input <-chan T) {
	if input == nil {
		panic("argument 'input' is mandatory")
	}

	for _, c := range c.subjects {
		c.Consume(input)
	}
}

func (s *MergeSubject[T]) Produce() <-chan T {
	var wg sync.WaitGroup
	output := make(chan T)

	outputFunc := func(c <-chan T) {
		defer wg.Done()
		for n := range c {
			select {
			case output <- n:
			case <-s.done:
				return
			}
		}
	}

	wg.Add(len(s.subjects))

	for _, s := range s.subjects {
		go outputFunc(s.Produce())
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}
