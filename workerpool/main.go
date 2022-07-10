package main

import (
	"fmt"
	"math/rand"
	"time"
)

var seed = rand.NewSource(time.Now().Unix())
var rng = rand.New(seed)

func main() {
	jobs := make(chan int, 10)
	res2 := make(chan string, 10)

	go worker(jobs, res2, "No1")
	go worker(jobs, res2, "No2")
	go worker(jobs, res2, "No3")

	for i := 0; i < 10; i++ {
		jobs <- i
	}
	close(jobs)

	for j := 0; j < 10; j++ {
		fmt.Println(<-res2)
	}
}

// worker function
func worker(jobs <-chan int, res chan<- string, id string) {
	for n := range jobs {
		msg := fmt.Sprintf("Worker Id: %s, Job Id: %d, res: %d", id, n, task(n+30))
		res <- msg
	}
}

// simulate time consuming task
func task(n int) int {
	x := rng.Intn(n)
	time.Sleep(time.Millisecond * time.Duration(x) * 100)
	return x
}
