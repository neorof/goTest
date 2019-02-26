package main

import (
	"fmt"
	"time"
)

var sig = make(chan int)
var sig2 = make(chan int)

func main() {
	fmt.Println(time.Now().Hour())
	go func() {
		go func() {
			sig2 <- 1
		}()
	}()

	select {
	case <-sig:
		fmt.Println("sig")
	case <-sig2:
		fmt.Println("sig2")
	}
}
