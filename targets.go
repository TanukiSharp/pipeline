package pipeline

import (
	"context"
	"sync"
)

// type inputInfo[T any] struct {
// 	channel     <-chan T
// 	linkContext context.Context
// 	unlink      UnlinkFunc
// }

// func newInputInfo[T any](input <-chan T) *inputInfo[T] {
// 	linkContext, unlink := context.WithCancel(context.Background())

// 	return &inputInfo[T]{
// 		channel:     input,
// 		linkContext: linkContext,
// 		unlink:      (UnlinkFunc)(unlink),
// 	}
// }

// ------------------------------------------------------------

type inputsSet[TInput, TOutput any] struct {
	sync.Mutex
	items map[<-chan TInput]*consumerInfo[TOutput]
}

type consumerInfo[TOutput any] struct {
	output chan TOutput
	unlink UnlinkFunc
}

func newConsumerInfo[TInput, TOutput any](ctx context.Context, input <-chan TInput, transformFunc func(input TInput) TOutput) *consumerInfo[TOutput] {
	output := make(chan TOutput)
	linkContext, unlink := context.WithCancel(context.Background())

	go func() {
		defer close(output)

		for {
			select {
			case n, ok := <-input:
				if ok == false {
					return
				}

				select {
				case output <- transformFunc(n):
				case <-ctx.Done():
					return
				case <-linkContext.Done():
					return
				}
			case <-ctx.Done():
				return
			case <-linkContext.Done():
				return
			}
		}
	}()

	return &consumerInfo[TOutput]{
		output: output,
		unlink: (UnlinkFunc)(unlink),
	}
}

func (inputs *inputsSet[TInput, TOutput]) set(ctx context.Context, input <-chan TInput, transformFunc func(input TInput) TOutput) *consumerInfo[TOutput] {
	inputs.Lock()
	defer inputs.Unlock()

	consumerInfo_, exists := inputs.items[input]

	if exists == false {
		consumerInfo_ = newConsumerInfo(ctx, input, transformFunc)
		inputs.items[input] = consumerInfo_
	}

	return consumerInfo_
}
