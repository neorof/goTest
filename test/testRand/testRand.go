package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// for i := 0; i < 10; i++ {
	// 	a := rand.Int()
	// 	fmt.Println(a)
	// }
	rand.Seed(time.Now().UnixNano()) //和时间相关的做种子不会重复
	for i := 0; i < 10; i++ {

		a := rand.Intn(100)
		fmt.Println(a)
	}
	// for i := 0; i < 10; i++ {
	// 	a := rand.Float32()
	// 	fmt.Println(a)
	// }
}
