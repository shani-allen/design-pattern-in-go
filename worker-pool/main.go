package main

import (
	"log"
	"runtime"
	"time"
)

//TODO: when to use worker pool
// limit the number of go routines or tasks that can be run concurrently and this way we manage the memory
// one use could be if we are having a web scraper  and we scraping the data from the multiple websites then
// there can be a case where our go out of memory while spawning the go routines for each website, to avoid
// this issue we can use worker pool to limit the number of go routines that can be run concurrently.

// we want to take generic type
type T = interface{}

// Worker is contract for worker pool inmplementation
type WorkerPool interface {
	Run()
	AddTask(task func())
}

// worker pool will hold maximum number of workers and queued task
type workerPool struct {
	maxWorkers int
	queuedTask chan func()
}

func NewWorkerPool(maxWorkers int) *workerPool {
	return &workerPool{
		maxWorkers: maxWorkers,
		queuedTask: make(chan func()),
	}
}

// Run will spawn the go routines based on the maxWorkers
func (w *workerPool) Run() {
	for i := 0; i < w.maxWorkers; i++ {
		go func(workerID int) {
			for task := range w.queuedTask {
				task()
			}
		}(i + 1)
	}
}

// AddTask will add the task to the queuedTask
func (w *workerPool) AddTask(task func()) {
	w.queuedTask <- task
}

func main() {
	//for monitoring purpose
	waitC := make(chan bool)
	go func() {
		for {
			log.Printf("[main] Total current goroutine: %d", runtime.NumGoroutine())
			time.Sleep(1 * time.Second)
		}
	}()
	totalWorkers := 5
	wp := NewWorkerPool(totalWorkers)
	wp.Run()

	type result struct {
		id    int
		value int
	}

	totalTask := 100
	resultC := make(chan result, totalTask)

	//writing the task to the channel
	for i := 0; i < totalTask; i++ {
		id := i + 1
		wp.AddTask(func() {
			log.Printf("[main] Starting task %d", id)
			time.Sleep(5 * time.Second)
			resultC <- result{id, id * 2}
		})
	}

	//reading from the channel to get the result
	for i := 0; i < totalTask; i++ {
		res := <-resultC
		log.Printf("[main] Task %d has been finished with result %d:", res.id, res.value)
	}

	<-waitC
}
