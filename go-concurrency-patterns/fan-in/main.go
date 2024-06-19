package main

import (
	"context"
	"fmt"
	"sync"
)

// basic fan-in pattern
// GO routines: a function that runs independently of the main function. We can think of it as process or lightweight thread.
// Channels: a way to communicate between go routines. there are two types of channels: buffered and unbuffered. buffered:
//can store multiple values, unbuffered: can store only one value.
// channels are blocking by default. meaning that always there should be some go routine to read from the channel.
// magic happens when we combine go routine and channels together, go routines communicate in a synchronized way and independent of each other.

// For Select done structure
// basically this pattern is to avoid any leak int the code.
// generally leak happens when any go routine is running and we do not need it. so to avoid this we need to send some signal that we do not need this go routine anymore.
// for this we use DONE channel

// PrintInt: this is the general pattern for select we use most of the places

func PrintInt(done <-chan struct{}, intStream <-chan int) {
	for {
		select {
		case res := <-intStream:
			fmt.Println(res)
		case <-done:
			return
		}
	}
}

// we can achieve the same effect by using the context also
func PrintIntContext(ctx context.Context, intStream <-chan int) {
	for {
		select {
		case res := <-intStream:
			fmt.Println(res)
		case <-ctx.Done():
			return
		}
	}
}

// use case for fan in pattern: lets say we are receiving the data from multiple sources and we want to combine them into single channel.
func FanIn(ctx context.Context, fetchers ...<-chan interface{}) <-chan interface{} {
	combineAll := make(chan interface{})

	var wg sync.WaitGroup
	wg.Add(len(fetchers))

	for _, f := range fetchers {
		f := f
		// spawning the go routines for each sources
		go func() {
			defer wg.Done()

			for {
				select {
				case res, ok := <-f:
					if !ok {
						return
					}
					combineAll <- res
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	defer func() {
		wg.Wait()
		close(combineAll)
	}()

	return combineAll
}
func main() {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})

	go func() {
		ch1 <- 1
		close(ch1)
	}()

	go func() {
		ch2 <- "hello"
		close(ch2)
	}()

	//ch2 := make(chan interface{})
	//ch2 <- "hello"
	combineAll := FanIn(context.Background(), ch1, ch2)

	for res := range combineAll {
		fmt.Println(res)
	}
}
