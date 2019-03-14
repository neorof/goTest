package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	people := People{"modishou", 18}
	buf,err := json.Marshal(people)
	if err!= nil {
		fmt.Println(err)
	}
	fmt.Println(string(buf))
}

type People struct{
	Name string
	Age int
}
