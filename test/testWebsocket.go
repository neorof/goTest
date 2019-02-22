package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kataras/golog"
	"net/url"
	"github.com/gorilla/websocket"
	"strconv"
	"os"
)

var addr = flag.String("addr", "localhost:9110", "ws service address")
//var addr = flag.String("addr", "172.31.20.65:9110", "ws service address")

const times  = 1
var cha1 = make(chan int)

func main() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/websocket"}
	file := newLogFile()
	golog.SetOutput(file)
	for i:=0; i<times; i++{
		doConnection(i, u)
	}
	val := <- cha1
	if val != 0 {
		golog.Error("id" + strconv.Itoa(val) + "error")
	}
}

func doConnection(id int, u url.URL) {
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(u.String(), nil)
	defer conn.Close()
	checkErr(err)
	for j := 0; j< 2; j++{
		if j == 0 {
			data := AuthData{"v1.0", 10, "06581ECE6C22300ADB41BF71" + strconv.Itoa(id)}
			msg := AuthRequest{"1_1", data}
			jsonMsg, err := json.Marshal(msg)
			checkErr(err)
			fmt.Println(string(jsonMsg))
			dataByte := []byte(string(jsonMsg))
			sendErr := conn.WriteMessage(websocket.TextMessage, dataByte)
			if sendErr != nil {
				cha1 <- 1
			} else {
				cha1 <- 0
			}
		} else {
			data := LogData{"aaa"}
			msg := LogRequest{"2_100", data}
			jsonMsg, err := json.Marshal(msg)
			checkErr(err)
			fmt.Println(string(jsonMsg))
			dataByte := []byte(string(jsonMsg))
			sendErr := conn.WriteMessage(websocket.TextMessage, dataByte)
			if sendErr != nil {
				cha1 <- 1
			}else {
				cha1 <- 0
			}
		}
	}
	//_, message, err := conn.ReadMessage()
	//fmt.Println(string(message))
	//checkErr(err)
}

func readMsg(conn websocket.Conn) {
	_, message, err := conn.ReadMessage()
	fmt.Println(string(message))
	checkErr(err)
}

func newLogFile() *os.File {
	filename := "benchTest0416_1.log"
	// open an output file, this will append to the today's file if server restarted.
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return f
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type AuthRequest struct {
	Cmd     string
	Payload AuthData
}

type AuthData struct {
	Xcx_version string
	Login_mode  int
	Token       string
}

type LogRequest struct {
	Cmd    string
	Paylod LogData
}

type LogData struct {
	Body string
}
