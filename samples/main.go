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

	source := pipeline.NewItemsSourceBlock(ctx, []int{0, 1, 2, 3, 4, 5}).SetName("source")
	propagator1 := pipeline.NewDelegatePropagatorBlock(ctx, func(input int) int {
		return input * input
	}).SetName("propagator1")
	propagator2 := pipeline.NewDelegatePropagatorBlock(ctx, func(input int) int {
		return input + 0
	}).SetName("propagator2")
	propagator3 := pipeline.NewDelegatePropagatorBlock(ctx, func(input int) any {
		fmt.Println(input)
		return nil
	}).SetName("propagator3")
	target := pipeline.NewDrainTargetBlock[any]().SetName("target")

	source.LinkTo(propagator1)
	propagator1.LinkTo(propagator2)
	propagator2.LinkTo(propagator3)
	propagator3.LinkTo(target)
}

func sample2() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	source := pipeline.NewItemsSourceBlock(ctx, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}).SetName("source")
	propagator1 := pipeline.NewDelegatePropagatorBlock(ctx, func(input int) int {
		return input * input
	}).SetName("propagator1")
	propagator2 := pipeline.NewDelegatePropagatorBlock(ctx, func(input int) int {
		if input > 25 {
			cancel()
		}
		return input + 0
	}).SetName("propagator2")
	target := pipeline.NewDelegateTargetBlock(func(input int) bool {
		fmt.Println(input)
		return true
	}).SetName("target")

	source.LinkTo(propagator1)
	propagator1.LinkTo(propagator2)
	propagator2.LinkTo(target)
}

func sample3() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	source := pipeline.NewItemsSourceBlock(ctx, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}).SetName("source")
	propagator1 := pipeline.NewDelegatePropagatorBlock(ctx, func(input int) int {
		return input * input
	}).SetName("propagator1")
	propagator2 := pipeline.NewDelegatePropagatorBlock(ctx, func(input int) int {
		return input + 0
	}).SetName("propagator2")
	target := pipeline.NewDelegateTargetBlock(func(input int) bool {
		fmt.Println(input)
		return input < 25
	}).SetName("target")

	source.LinkTo(propagator1)
	propagator1.LinkTo(propagator2)
	propagator2.LinkTo(target)

	// Because the target consumer has exited early,
	// we need to cancel the pipeline to release goroutines of upstream layers.
	// This is an explicit cause because it is a specific case, but the defer cancel anyways.
	cancel()
}

func sample4() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	source := pipeline.NewItemsSourceBlock(ctx, []int{4, 9}).SetName("source")

	propagator1 := pipeline.NewDelegatePropagatorBlock(ctx, func(input int) int {
		return input * input
	}).SetName("propagator1")
	propagator2 := pipeline.NewDelegatePropagatorBlock(ctx, func(input int) int {
		return input * input
	}).SetName("propagator2")
	merge := pipeline.NewMergePropagatorBlockWithInstances[int, int](ctx, propagator1, propagator2).SetName("merge")

	target := pipeline.NewDelegateTargetBlock(func(input int) bool {
		fmt.Println(input)
		return true
	}).SetName("target")

	source.LinkTo(merge)
	merge.LinkTo(target)
}

type sample5Tuple struct {
	in  int
	out bool
}

func sample5() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	source := pipeline.NewItemsSourceBlock(ctx, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}).SetName("source")
	propagator1 := pipeline.NewDelegatePropagatorBlock(ctx, func(input int) int {
		return input * input
	}).SetName("propagator1")
	propagator2 := pipeline.NewDelegatePropagatorBlock(ctx, func(input int) int {
		return input + 0
	}).SetName("propagator2")
	propagator3 := pipeline.NewDelegatePropagatorBlock(ctx, func(input int) sample5Tuple {
		return sample5Tuple{in: input, out: input < 25}
	}).SetName("propagator3")

	source.LinkTo(propagator1)
	propagator1.LinkTo(propagator2)
	propagator2.LinkTo(propagator3)

	// Consume manually...
	ch := propagator3.Produce()

	// ...all with a loop.
	for n := range ch {
		fmt.Printf("i: %d -> n: %v\n", n.in, n.out)
	}

	// ...or some but not all, for testing.
	// print := func(value sample5Tuple) {
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
