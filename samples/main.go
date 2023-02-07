package main

import (
	"context"
	"fmt"
	"pipeline"
)

// ================================================================================

func sample1() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	source := pipeline.NewItemsProducer(ctx, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	subject1 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input * input
	})
	subject2 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input + 0
	})
	subject3 := pipeline.NewDelegateSubject(ctx, func(input int) any {
		fmt.Println(input)
		return nil
	})
	sink := pipeline.NewDrainConsumer[any]()

	source.LinkTo(subject1)
	subject1.LinkTo(subject2)
	subject2.LinkTo(subject3)
	subject3.LinkTo(sink)

	// TODO: Something to do here :/
}

func sample2() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	source := pipeline.NewItemsProducer(ctx, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	subject1 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input * input
	})
	subject2 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		if input > 25 {
			cancel()
		}
		return input + 0
	})
	sink := pipeline.NewDelegateConsumer(func(input int) bool {
		fmt.Println(input)
		return true
	})

	source.LinkTo(subject1)
	subject1.LinkTo(subject2)
	subject2.LinkTo(sink)

	// TODO: Something to do here :/
}

func sample3() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	source := pipeline.NewItemsProducer(ctx, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	subject1 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input * input
	})
	subject2 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input + 0
	})
	sink := pipeline.NewDelegateConsumer(func(input int) bool {
		fmt.Println(input)
		return input < 25
	})

	source.LinkTo(subject1)
	subject1.LinkTo(subject2)
	subject2.LinkTo(sink)

	// TODO: Something to do here :/

	// Because the sink consumer has exited early,
	// we need to cancel the pipeline to release goroutines of upstream layers.
	// This is an explicit cause because it is a specific case, but the defer cancel anyways.
	cancel()
}

func sample4() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	source := pipeline.NewItemsProducer(ctx, []int{4, 9})

	subject1 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input * input
	})
	subject2 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input * input
	})
	merge := pipeline.NewMergeSubjectWithInstances[int, int](ctx, subject1, subject2)

	sink := pipeline.NewDelegateConsumer(func(input int) bool {
		fmt.Println(input)
		return true
	})

	source.LinkTo(merge)
	merge.LinkTo(sink)

	// TODO: Something to do here :/
}

func runSample(name string, f func()) {
	fmt.Printf("--- %s -------------------\n", name)
	f()
	fmt.Println("-------------------------------")
	fmt.Println()
}

func main() {
	runSample("sample1", sample1)
	runSample("sample2", sample2)
	runSample("sample3", sample3)
	runSample("sample4", sample4)
}
