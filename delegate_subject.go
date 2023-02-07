package pipeline

import "context"

type DelegateSubject[TInput any, TOutput any] struct {
	ctx           context.Context
	output        chan TOutput
	transformFunc func(input TInput) TOutput
}

var _ Subject[any, any] = &DelegateSubject[any, any]{}

func NewDelegateSubject[TInput any, TOutput any](ctx context.Context, transformFunc func(input TInput) TOutput) *DelegateSubject[TInput, TOutput] {
	if ctx == nil {
		panic("argument 'ctx' is mandatory")
	}
	if transformFunc == nil {
		panic("argument 'transformFunc' is mandatory")
	}

	return &DelegateSubject[TInput, TOutput]{
		ctx:           ctx,
		transformFunc: transformFunc,
	}
}

func (s *DelegateSubject[TInput, TOutput]) Consume(input <-chan TInput) func() {
	if input == nil {
		panic("argument 'input' is mandatory")
	}

	brk := false
	s.output = make(chan TOutput)

	go func() {
		defer close(s.output)
		for n := range input {
			select {
			case s.output <- s.transformFunc(n):
			case <-s.ctx.Done():
				return
			default:
				if brk {
					break
				}
			}
		}
	}()

	return func() {
		brk = true
	}
}

func (s *DelegateSubject[TTInput, TOutput]) Produce() <-chan TOutput {
	return s.output
}

func (s *DelegateSubject[TInput, TOutput]) LinkTo(consumer Consumer[TOutput]) func() {
	return consumer.Consume(s.Produce())
}
