package main

import (
	"fmt"
	"os"
	"os/signal"
)
//
func main() {

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, os.Kill)
	s := <- signalCh

	fmt.Println("中断信号:", s)
}
