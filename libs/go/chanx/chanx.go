// Package chanx implements various concurrency patterns
// inspired by
// www.oreilly.com/library/view/concurrency-in-go/9781491941294/ch04.html
// with generics
//
//revive:disable:cognitive-complexity
package chanx

import (
	"context"
	"sync"
)

// OrDone either forwards the values of c, or quits if ctx is done.
func OrDone[T any](ctx context.Context, c <-chan T) <-chan T {
	values := make(chan T)
	go func() {
		defer close(values)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-c:
				if !ok {
					return
				}

				select {
				case values <- v:
				case <-ctx.Done():
				}
			}
		}
	}()

	return values
}

// Tee receives T values from in and forwards them to left and right channels
func Tee[T any](ctx context.Context, in <-chan T) (_, _ <-chan T) {
	left := make(chan T)
	right := make(chan T)
	go func() {
		defer close(left)
		defer close(right)

		for val := range OrDone(ctx, in) {
			var left, right = left, right

			for i := 0; i < 2; i++ {
				select {
				case <-ctx.Done():
				case left <- val:
					left = nil
				case right <- val:
					right = nil
				}
			}
		}
	}()

	return left, right
}

// FanIn drains in parallel a number of channels [multiplex].
// Returns the values in a single channel
func FanIn[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	wg.Add(len(channels))

	multiplexed := make(chan T)

	drain := func(c <-chan T) {
		defer wg.Done()
		for val := range c {
			select {
			case <-ctx.Done():
				return
			case multiplexed <- val:
			}
		}
	}

	for _, c := range channels {
		go drain(c)
	}

	go func() {
		wg.Wait()
		close(multiplexed)
	}()

	return multiplexed
}

// Drain drains the content of c into multiplexed.
// This func is supposed to call N times to Drain N channels
// into a single one
func Drain[T any](
	ctx context.Context,
	c <-chan T,
	wg *sync.WaitGroup,
	multiplexed chan<- T,
) {
	defer wg.Done()
	for val := range c {
		select {
		case <-ctx.Done():
			return
		case multiplexed <- val:
		}
	}
}
