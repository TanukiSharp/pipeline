package pipeline

import (
	"testing"
)

func TestRegularDelegateConsumer(t *testing.T) {
	obtained := 0

	consumer := NewDelegateConsumer(func(input int) bool {
		obtained += input
		return true
	})

	ch := make(chan int)

	go func() {
		defer close(ch)
		ch <- 1
		ch <- 2
		ch <- 3
	}()

	consumer.Consume(ch)

	const expected = 6

	if obtained != expected {
		t.Fatalf("expected %d, obtained %d\n", expected, obtained)
	}
}

func TestInterruptedDelegateConsumer(t *testing.T) {
	obtained := 0

	consumer := NewDelegateConsumer(func(input int) bool {
		if input > 4 {
			return false
		}
		obtained += input
		return true
	})

	ch := make(chan int)

	go func() {
		defer close(ch)
		ch <- 1
		ch <- 2
		ch <- 3
		ch <- 4
		ch <- 5
		ch <- 6
		ch <- 7
		ch <- 8
	}()

	consumer.Consume(ch)

	const expected = 10

	if obtained != expected {
		t.Fatalf("expected %d, obtained %d\n", expected, obtained)
	}
}
