package main

import (
	"fmt"
	"pipeline"
)

// ================================================================================

func sample1() {
	p := pipeline.NewPipeline[int]()

	source := pipeline.NewItemsProducer(p, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	subject1 := pipeline.NewDelegateSubject(p, func(input int) int {
		return input * input
	})
	subject2 := pipeline.NewDelegateSubject(p, func(input int) int {
		return input + 0
	})
	sink := pipeline.NewDelegateConsumer(func(input int) bool {
		fmt.Println(input)
		return true
	})

	p.RegisterSource(source)
	p.AddSubject(subject1)
	p.AddSubject(subject2)
	p.RegisterSink(sink)

	err := p.Run()

	if err != nil {
		panic(err)
	}
}

func sample2() {
	p := pipeline.NewPipeline[int]()

	source := pipeline.NewItemsProducer(p, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	subject1 := pipeline.NewDelegateSubject(p, func(input int) int {
		return input * input
	})
	subject2 := pipeline.NewDelegateSubject(p, func(input int) int {
		if input > 25 {
			p.Cancel()
		}
		return input + 0
	})
	sink := pipeline.NewDelegateConsumer(func(input int) bool {
		fmt.Println(input)
		return true
	})

	p.RegisterSource(source)
	p.AddSubject(subject1)
	p.AddSubject(subject2)
	p.RegisterSink(sink)

	err := p.Run()

	if err != nil {
		panic(err)
	}
}

func sample3() {
	p := pipeline.NewPipeline[int]()

	source := pipeline.NewItemsProducer(p, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	subject1 := pipeline.NewDelegateSubject(p, func(input int) int {
		return input * input
	})
	subject2 := pipeline.NewDelegateSubject(p, func(input int) int {
		return input + 0
	})
	sink := pipeline.NewDelegateConsumer(func(input int) bool {
		fmt.Println(input)
		return input < 25
	})

	p.RegisterSource(source)
	p.AddSubject(subject1)
	p.AddSubject(subject2)
	p.RegisterSink(sink)

	err := p.Run()

	if err != nil {
		panic(err)
	}

	// Because the sink consumer has exited early,
	// we need to cancel the pipeline to release goroutines of upstream layers.
	p.Cancel()
}

func sample4() {
	p := pipeline.NewPipeline[int]()

	source := pipeline.NewItemsProducer(p, []int{4, 9})

	subject1 := pipeline.NewDelegateSubject(p, func(input int) int {
		return input * input
	})
	subject2 := pipeline.NewDelegateSubject(p, func(input int) int {
		return input * input
	})
	merge := pipeline.NewMergeSubjectWithInstances[int](p, subject1, subject2)

	sink := pipeline.NewDelegateConsumer(func(input int) bool {
		fmt.Println(input)
		return true
	})

	p.RegisterSource(source)
	p.AddSubject(merge)
	p.RegisterSink(sink)

	err := p.Run()

	if err != nil {
		panic(err)
	}
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
