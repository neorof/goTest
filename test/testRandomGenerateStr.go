package main



import (
	"fmt"
	"math/rand"
	"time"
)

const str = "0123456789abcdefghijklmnopqrstuvwxyz"

func main() {
	fmt.Println(GetRandomToken())
}

func GetRandomToken() string{
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}