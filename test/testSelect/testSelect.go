package main

import (
	"fmt"
	"runtime"
)
//在执行select语句的时候，运行时系统会自上而下地判断每个case中的发送或接收操作是否可以被立即执行(立即执行：意思是当前Goroutine不会因此操作而被阻塞)
func main(){
	runtime.GOMAXPROCS(1)
	int_chan := make(chan int, 1)
	str_chan := make(chan string, 1)
	int_chan <- 1
	str_chan <- "hello"
	select {
	case val := <- int_chan:
		fmt.Println(val)
	case val := <- str_chan:
		fmt.Println(val)
	}
	fmt.Println(123)
}