package main

import (
	"flag"
	"fmt"
)

func main() {
	conn := flag.Int("conn", -1, "Input conn")
	if *conn == -1 {
		fmt.Println("please input conn param!")
		return
	}
	flag.Parse()

	fmt.Println(3 + *conn)
}
