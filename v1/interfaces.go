package pipeline

type Producer[T any] interface {
	Produce() <-chan T
}

type Consumer[T any] interface {
	Consume(<-chan T)
}

type Subject[T any] interface {
	Producer[T]
	Consumer[T]
}
