package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/golang/snappy"
)

//测试readMessage是会被阻塞
const uri = "ws://localhost:9110/websocket"

func main() {

	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	var data = "{\"cmd\":\"1_1\",\"payload\":{\"login_mode\":11,\"token\":\"aaaaaa\",\"xcx_version\":\"v10.1\"}}"
	byteData := snappy.Encode(nil, []byte(data))
	err1 := conn.WriteMessage(websocket.BinaryMessage , byteData)
	if err1 != nil {
		panic(err1)
	}

	_, recv, _ := conn.ReadMessage()
	uncompress,_ := snappy.Decode(nil, recv)
	fmt.Println(string(uncompress))
	// fmt.Println(string(recv))

}
