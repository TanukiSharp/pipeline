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

	source := pipeline.NewItemsProducer(ctx, []int{0, 1, 2, 3, 4, 5}).SetName("source")
	subject1 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input * input
	}).SetName("subject1")
	subject2 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input + 0
	}).SetName("subject2")
	subject3 := pipeline.NewDelegateSubject(ctx, func(input int) any {
		fmt.Println(input)
		return nil
	}).SetName("subject3")
	sink := pipeline.NewDrainConsumer[any]().SetName("sink")

	source.LinkTo(subject1)
	subject1.LinkTo(subject2)
	subject2.LinkTo(subject3)
	subject3.LinkTo(sink)
}

func sample2() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	source := pipeline.NewItemsProducer(ctx, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}).SetName("source")
	subject1 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input * input
	}).SetName("subject1")
	subject2 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		if input > 25 {
			cancel()
		}
		return input + 0
	}).SetName("subject2")
	sink := pipeline.NewDelegateConsumer(func(input int) bool {
		fmt.Println(input)
		return true
	}).SetName("sink")

	source.LinkTo(subject1)
	subject1.LinkTo(subject2)
	subject2.LinkTo(sink)
}

func sample3() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	source := pipeline.NewItemsProducer(ctx, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}).SetName("source")
	subject1 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input * input
	}).SetName("subject1")
	subject2 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input + 0
	}).SetName("subject2")
	sink := pipeline.NewDelegateConsumer(func(input int) bool {
		fmt.Println(input)
		return input < 25
	}).SetName("sink")

	source.LinkTo(subject1)
	subject1.LinkTo(subject2)
	subject2.LinkTo(sink)

	// Because the sink consumer has exited early,
	// we need to cancel the pipeline to release goroutines of upstream layers.
	// This is an explicit cause because it is a specific case, but the defer cancel anyways.
	cancel()
}

func sample4() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	source := pipeline.NewItemsProducer(ctx, []int{4, 9}).SetName("source")

	subject1 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input * input
	}).SetName("subject1")
	subject2 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input * input
	}).SetName("subject2")
	merge := pipeline.NewMergeSubjectWithInstances[int, int](ctx, subject1, subject2).SetName("merge")

	sink := pipeline.NewDelegateConsumer(func(input int) bool {
		fmt.Println(input)
		return true
	}).SetName("sink")

	source.LinkTo(merge)
	merge.LinkTo(sink)
}

func sample5() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	source := pipeline.NewItemsProducer(ctx, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}).SetName("source")
	subject1 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input * input
	}).SetName("subject1")
	subject2 := pipeline.NewDelegateSubject(ctx, func(input int) int {
		return input + 0
	}).SetName("subject2")
	subject3 := pipeline.NewDelegateSubject(ctx, func(input int) struct {
		in  int
		out bool
	} {
		return struct {
			in  int
			out bool
		}{in: input, out: input < 25}
	}).SetName("subject3")

	source.LinkTo(subject1)
	subject1.LinkTo(subject2)
	subject2.LinkTo(subject3)

	// Consume manually...
	ch := subject3.Produce()

	// ...all with a loop.
	for n := range ch {
		fmt.Printf("i: %d -> n: %v\n", n.in, n.out)
	}

	// ...or some but not all, for testing.
	// print := func(value struct {
	// 	in  int
	// 	out bool
	// }) {
	// 	fmt.Printf("i: %d -> n: %v\n", value.in, value.out)
	// }
	// print(<-ch)
	// print(<-ch)
	// print(<-ch)
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
	runSample("sample5", sample5)
}
