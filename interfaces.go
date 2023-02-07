package pipeline

type Producer[TOutput any] interface {
	Produce() <-chan TOutput
	LinkTo(consumer Consumer[TOutput]) func()
}

type Consumer[TInput any] interface {
	Consume(<-chan TInput) func()
}

type Subject[TInput any, TOutput any] interface {
	Consumer[TInput]
	Producer[TOutput]
}
