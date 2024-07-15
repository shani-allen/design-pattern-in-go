package main

//Replace time.Sleep with channel communication

import (
	"fmt"
)

func main() {
	ch := make(chan int)

	//var wg sync.WaitGroup
	//wg.Add(1)
	go func(<-chan int) {
		//defer wg.Done()
		fmt.Println("hello from goroutine!")
		ch <- 1

	}(ch)

	fmt.Println(<-ch)
	//wg.Wait()
	//time.Sleep(time.Second * 1)
}
