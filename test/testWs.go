package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)
//测试readMessage是会被阻塞
const uri = "ws://172.31.20.65:9110/websocket"

func main() {

	conn,_, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	var data = "{\"cmd\":\"1_1\",\"ctx\":1,\"payload\":{\"xcx_version\":\"v3.1.1\",\"login_mode\":10,\"token\":\"123456\"}}"
	err1 := conn.WriteMessage(websocket.TextMessage, []byte(data))
	if err1 != nil {
		panic(err1)
	}

	_, recv, _ := conn.ReadMessage()
	fmt.Println(string(recv))

}
