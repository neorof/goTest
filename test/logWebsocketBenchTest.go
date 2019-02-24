package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kataras/golog"
)

const url = "ws://172.31.20.65:9110/websocket"

//const url = "ws://127.0.0.1:9110/websocket"

var flagToWaitChan = make(chan int)

func main() {
	conn := flag.Int("conn", -1, "Input conn")
	flag.Parse()
	if *conn == -1 {
		fmt.Println("please input conn param!")
		return
	}
	f := newLogFile()
	defer f.Close()
	golog.SetOutput(f)
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	var timeToSleep = 5
	for curConn := 0; curConn < *conn; curConn++ {
		golog.Info("new conn id" + strconv.Itoa(curConn) + "_")
		go newOneConn(curConn)
		// 1000以下sleep 0.3
		var flag = <-flagToWaitChan
		if flag == 1 {
			golog.Info("id" + strconv.Itoa(curConn) + "_ sleep " + strconv.Itoa(timeToSleep))
			time.Sleep(time.Duration(timeToSleep))
			timeToSleep = rand.Intn(60) + 10
		}
	}
	select {
	case <-interrupt:
		fmt.Println("人工中断")
	}
	//time.Sleep(time.Duration(time.Hour)) //防止主线程结束中断其他线程
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

func newOneConn(id int) {

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		golog.Error("id"+strconv.Itoa(id)+"_ dial error:", err)
		flagToWaitChan <- 1
		return
	}
	golog.Info("id" + strconv.Itoa(id) + "_ dial success")
	flagToWaitChan <- 0
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		var errTime = 0
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				golog.Error("id"+strconv.Itoa(id)+"_ read error:", err)
				errTime++
				time.Sleep(time.Duration(4*errTime) * time.Second)
				golog.Warn("id" + strconv.Itoa(id) + "_ read error times: " + strconv.Itoa(errTime))
				if errTime > 10 {
					return
				}
				continue
			}
			//log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var errTime2 = 0
	var reqTime uint64 = 0
	var jsonMsg []byte
	for {
		reqTime++
		select {
		case <-done:
			return
		case <-ticker.C:
			if reqTime == 1 {
				data := AuthData{"v1.0", 10, " 023ZJLl12DYkOU0DFSm12dsLl12ZJLl" + strconv.Itoa(id)}
				msg := AuthRequest{"1_1", data}
				jsonMsg, _ = json.Marshal(msg)
			} else {
				data := LogData{"aaa"}
				msg := LogRequest{"2_100", data}
				jsonMsg, _ = json.Marshal(msg)
			}

			err = c.WriteMessage(websocket.TextMessage, jsonMsg)
			if err != nil {
				golog.Error("id"+strconv.Itoa(id)+"_ Write error:", err)
				errTime2++
				time.Sleep(time.Duration(4*errTime2) * time.Second)
				if errTime2 > 10 {
					return
				}
				continue
			}

		case <-interrupt:
			golog.Debug("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				golog.Error("id"+strconv.Itoa(id)+"_ Close write error:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func putFlag() {
	flagToWaitChan <- 0
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
