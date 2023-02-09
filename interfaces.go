package pipeline

type UnlinkFunc func()

type SourceBlock[TOutput any] interface {
	Produce() <-chan TOutput
	LinkTo(target TargetBlock[TOutput]) UnlinkFunc
}

type TargetBlock[TInput any] interface {
	Consume(<-chan TInput) UnlinkFunc
}

type PropagatorBlock[TInput, TOutput any] interface {
	SourceBlock[TOutput]
	TargetBlock[TInput]
}
