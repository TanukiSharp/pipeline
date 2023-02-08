package pipeline

type UnlinkFunc func()

type Producer[TOutput any] interface {
	Produce() <-chan TOutput
	LinkTo(consumer Consumer[TOutput]) UnlinkFunc
}

type Consumer[TInput any] interface {
	Consume(<-chan TInput) UnlinkFunc
}

type Subject[TInput, TOutput any] interface {
	Consumer[TInput]
	Producer[TOutput]
}
