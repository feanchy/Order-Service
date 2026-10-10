package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	jobs := make(chan int)

	const countWorkers = 5

	wg.Add(countWorkers)
	for i := 0; i < countWorkers; i++ {
		go worker(&wg, i, jobs)
	}

	for i := 0; i < 20; i++ {
		jobs <- i
	}

	close(jobs)

	wg.Wait()

	fmt.Println("all jobs is completed")
}

func worker(wg *sync.WaitGroup, workerID int, jobID <-chan int) {
	defer wg.Done()

	for job := range jobID {
		fmt.Printf("worker %d, start job %d\n", workerID, job)
		time.Sleep(1 * time.Second)
		fmt.Printf("worker %d is done job %d\n", workerID, job)
	}

}
