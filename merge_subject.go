package pipeline

import (
	"context"
	"sync"
)

type MergeSubject[TInput, TOutput any] struct {
	name     string
	ctx      context.Context
	subjects []Subject[TInput, TOutput]
}

var _ Subject[any, any] = &MergeSubject[any, any]{}

func NewMergeSubjectWithFactory[TInput, TOutput any](ctx context.Context, howMany int, factoryFunc func(i int) Subject[TInput, TOutput]) *MergeSubject[TInput, TOutput] {
	if ctx == nil {
		panic("argument 'ctx' is mandatory")
	}
	if factoryFunc == nil {
		panic("argument 'factoryFunc' is mandatory")
	}

	subjects := make([]Subject[TInput, TOutput], howMany)

	for i := 0; i < howMany; i++ {
		subjects[i] = factoryFunc(i)
	}

	return NewMergeSubjectWithInstances(ctx, subjects...)
}

func NewMergeSubjectWithInstances[TInput, TOutput any](ctx context.Context, subjects ...Subject[TInput, TOutput]) *MergeSubject[TInput, TOutput] {
	if ctx == nil {
		panic("argument 'ctx' is mandatory")
	}
	if subjects == nil {
		panic("argument 'subject' is mandatory")
	}
	if len(subjects) == 0 {
		panic("argument 'subjects' must contain at least one element")
	}

	return &MergeSubject[TInput, TOutput]{
		ctx:      ctx,
		subjects: subjects,
	}
}

func (c *MergeSubject[TInput, TOutput]) GetName() string {
	return c.name
}

func (c *MergeSubject[TInput, TOutput]) SetName(name string) *MergeSubject[TInput, TOutput] {
	c.name = name
	return c
}

func (c *MergeSubject[TInput, TOutput]) Consume(input <-chan TInput) UnlinkFunc {
	if input == nil {
		panic("argument 'input' is mandatory")
	}

	unlinkFuncs := []func(){}

	for _, c := range c.subjects {
		unlinkFunc := c.Consume(input)
		unlinkFuncs = append(unlinkFuncs, unlinkFunc)
	}

	return func() {
		for _, unlinkFunc := range unlinkFuncs {
			unlinkFunc()
		}
	}
}

func (s *MergeSubject[TInput, TOutput]) Produce() <-chan TOutput {
	var wg sync.WaitGroup
	output := make(chan TOutput)

	outputFunc := func(c <-chan TOutput) {
		defer wg.Done()
		for n := range c {
			select {
			case output <- n:
			case <-s.ctx.Done():
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

func (s *MergeSubject[TInput, TOutput]) LinkTo(consumer Consumer[TOutput]) UnlinkFunc {
	return consumer.Consume(s.Produce())
}
