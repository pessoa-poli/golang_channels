package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const numOfWorks = 200
const maxWorkers = 20

//var mutex = sync.Mutex{}

var listOfStuff = []int{}

func grabItemFromList() []int {
	//mutex.Lock()
	if len(listOfStuff) == 0 {
		return []int{}
	}
	item := listOfStuff[len(listOfStuff)-1]
	listOfStuff = listOfStuff[:len(listOfStuff)-1]
	//mutex.Unlock()
	return []int{item}
}

func fillList(numOfItems int) {
	for i := 0; i < numOfItems; i++ {
		//listOfStuff = append(listOfStuff, rand.Intn(1000))
		listOfStuff = append(listOfStuff, i)
	}
}

func doSomeWork(work int, c chan<- int) {
	defer func() { c <- work }()
	fmt.Printf("Doing work %d\n", work)
	workThisMuch, _ := time.ParseDuration(fmt.Sprintf("%vs", rand.Intn(5-1)+1))
	//fmt.Printf("Going to work %v \n", workThisMuch)
	time.Sleep(workThisMuch)
}

var totalWG = sync.WaitGroup{}

func manager(c chan int) {
	worksleft := numOfWorks
	for {
		//Wait till a worker has finished
		workDone := <-c
		worksleft--
		fmt.Printf("Finished work %v\n", workDone)
		if worksleft == 0 {
			fmt.Println("Finished all work.")
			totalWG.Done()
		}
		//Check if there is work left to do
		if maxWorkers < numOfWorks {
			nextItem := grabItemFromList()
			if len(nextItem) == 0 {
				continue
			}
			fmt.Println("Sending new worker.")
			go doSomeWork(nextItem[0], c)
		}
	}

}

func main() {
	fmt.Println("Hello, playground")
	fillList(numOfWorks)
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	c := make(chan int)

	for i := 0; i < maxWorkers; i++ {
		item := grabItemFromList()
		if len(item) == 0 {
			break
		}
		go doSomeWork(item[0], c)
	}
	fmt.Println("Launching manager.")
	time.Sleep(8 * time.Second)
	go manager(c)
	totalWG.Add(1)
	totalWG.Wait()
	fmt.Println("Work Done")
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
