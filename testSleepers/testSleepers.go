package main

import (
	"fmt"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func sleeper1() {
	time.Sleep(10 * time.Second)
	wg.Done()
}

func sleeper2() {
	time.Sleep(12 * time.Second)
	wg.Done()
}

func main() {
	fmt.Println("Hello")
	go sleeper1()
	go sleeper2()
	wg.Add(2)
	wg.Wait()
	fmt.Println("Finished")

}
