package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func doSomeWork(work int, c chan<- int) {
	defer func() { c <- work }()
	rand.Seed(time.Now().UnixNano())
	fmt.Printf("Doing work %d\n", work)
	workThisMuch, _ := time.ParseDuration(fmt.Sprintf("%vs", rand.Intn(5-1)+1))
	//fmt.Printf("Going to work %v \n", workThisMuch)
	time.Sleep(workThisMuch)
}

var totalWG = sync.WaitGroup{}

func manager(c chan int, totalWorks, maxWorkers int) {
	worksLeftToDo := totalWorks
	nextWorkToDispatch := maxWorkers
	for worksLeftToDo >= 0 {
		workFinished := <-c
		worksLeftToDo--
		fmt.Printf("Work %v finished.\n", workFinished)
		if totalWorks > maxWorkers && nextWorkToDispatch != totalWorks {
			go doSomeWork(nextWorkToDispatch, c)
			nextWorkToDispatch++
		}
		if worksLeftToDo == 0 {
			break
		}
	}
	totalWG.Done()
}

func main() {
	fmt.Println("Hello, playground")
	start := time.Now()
	c := make(chan int)
	totalWorks := 200
	maxWorkers := 20
	for i := 0; i < maxWorkers; i++ {
		go doSomeWork(i, c)
	}
	go manager(c, totalWorks, maxWorkers)
	totalWG.Add(1)
	totalWG.Wait()
	fmt.Println("Work Done")
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
