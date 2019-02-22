package main

import "fmt"

var ch = make(chan int)

func main() {
	go beforeDefer()
	<- ch
	fmt.Println("main done")
}

func putChan() {
	ch <- 1
}

func beforeDefer() {
	defer putChan()
	defer fmt.Println("cccc")
	defer fmt.Println("aaaa")
	fmt.Println("bbbb")
}
