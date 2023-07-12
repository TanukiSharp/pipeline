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

type sample4Tuple struct {
	original int
	output   int
}

func sample4() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sourceItems := []sample4Tuple{
		{original: 0},
		{original: 1},
		{original: 2},
		{original: 3},
		{original: 4},
		{original: 5},
	}

	source := pipeline.NewItemsSourceBlock(ctx, sourceItems).SetName("source")

	p1Counter := 0
	p2Counter := 0
	targetCounter := 0

	propagator1 := pipeline.NewDelegatePropagatorBlock(ctx, func(input sample4Tuple) sample4Tuple {
		p1Counter++
		return sample4Tuple{original: input.original, output: input.original * input.original}
	}).SetName("propagator1")
	propagator2 := pipeline.NewDelegatePropagatorBlock(ctx, func(input sample4Tuple) sample4Tuple {
		p2Counter++
		return sample4Tuple{original: input.original, output: input.original * input.original}
	}).SetName("propagator2")

	target := pipeline.NewDelegateTargetBlock(func(input sample4Tuple) bool {
		targetCounter++
		fmt.Printf("original: %d, output: %d\n", input.original, input.output)
		return true
	}).SetName("target")

	// merge := pipeline.NewMergePropagatorBlockWithInstances[sample4Tuple, sample4Tuple](ctx, propagator1, propagator2).SetName("merge")
	// source.LinkTo(merge)
	// merge.LinkTo(target)

	source.LinkTo(propagator1)
	source.LinkTo(propagator2)
	propagator1.LinkTo(target)
	propagator2.LinkTo(target)

	pCounter := p1Counter + p2Counter

	if len(sourceItems) != pCounter || pCounter != targetCounter {
		panic("Something went wrong.")
	}
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

func sample6() {
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

func runSample(name string, f func()) {
	fmt.Printf("--- %s -------------------\n", name)
	f()
	fmt.Println("-------------------------------")
	fmt.Println()
}

func main() {
	// runSample("sample1", sample1)
	// runSample("sample2", sample2)
	// runSample("sample3", sample3)
	runSample("sample4", sample4)
	// runSample("sample5", sample5)
}
