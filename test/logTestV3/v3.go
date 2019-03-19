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

// const uri = "ws://172.31.20.65:9110/websocket"

const uri = "ws://127.0.0.1:9110/websocket"
const str = "0123456789abcdefghijklmnopqrstuvwxyz"

const logDa1 = "wxapp_type=xcx&platform=wxapp&app=pdd&page_name=index&page_id=10002_1551430328032_8V4bP3wwBL&page_sn=10002&page_url=pages%2Findex%2Findex&page=index&op=epv&sub_op=leave&page_duration=66345&enter_time=1551430328084&user_id=6677178664&log_id=15514303944293t0jBYCVdLCvQHqd&token=LKIG7BPWHL7TPPEORCOOEG6SCG3JPCTLG74CIXJPM4LII5X7NJ3Q101287d&app_version=v3.1.6.2&time=1551430394299&xcx_x_scene=1089&pdd_user_type=0&xcx_trace_id=bQH0A18KIhqkBqZSvBjTDnscsLTZPjBB&session_id=8NLzGS5Vuq9eNmCLvj82ZlhuvdegJoC7&greyscale_guid=cf1471a9a152de6fe27e4d23423c7eda&withsocket=1&openId="

var readErr = make(chan int)
var writeErr = make(chan int)
var openErr = make(chan int)
var mo int
var connTime int

func main() {
	conn := flag.Int("c", -1, "Input conn")
	mode := flag.Int("m", -1, "Input mode")
	flag.Parse()
	if *conn == -1 || *mode == -1 {
		fmt.Println("please input param(-c -m)!")
		return
	}
	// f := newLogFile()
	// defer f.Close()
	// golog.SetOutput(f)
	mo = *mode
	fmt.Println("auth mode=", *mode)

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	for curConn := 0; curConn < *conn; curConn++ {
		time.Sleep(20 * time.Microsecond)
		go newOneConn(curConn)
	}

	select {
	case <-interrupt:
		fmt.Println("人工中断")
	case <-openErr:
		fmt.Println("建立链接异常")
	case <-readErr:
		fmt.Println("读取异常")
	case <-writeErr:
		fmt.Println("发送异常")
	}
	//time.Sleep(time.Duration(time.Hour)) //防止主线程结束中断其他线程
}

func newOneConn(id int) {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	c, _, err := websocket.DefaultDialer.Dial(uri, nil)

	if err != nil {
		fmt.Println("id:" + strconv.Itoa(id) + " connect fail")
		connTime++
		if connTime >= 2 {
			openErr <- 1
			return
		}
	}

	defer c.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		var errTime = 0
		for {
			if c != nil {
				_, _, err := c.ReadMessage()
				if err != nil {
					errTime++
					if errTime >= 2 {
						readErr <- 1
						return
					}
					continue
				}
			}
		}
	}()

	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()

	var errTime2 int
	var reqTime uint64
	var jsonMsg []byte
	var openID string
	for {
		reqTime++
		select {
		case <-done:
			return
		case <-ticker.C:
			if reqTime == 1 {
				openID = "mo" + GetRandomToken() + strconv.Itoa(id)
				data := AuthData{"v1.0", mo, openID}
				msg := AuthRequest{"1_1", data}
				jsonMsg, _ = json.Marshal(msg)
			} else {
				randomStr := logDa1 + openID
				data := LogData{randomStr}
				msg := LogRequest{"2_101", data}
				jsonMsg, _ = json.Marshal(msg)
			}

			if c != nil {
				err = c.WriteMessage(websocket.TextMessage, jsonMsg)
				if err != nil {
					errTime2++
					if errTime2 >= 2 {
						writeErr <- 1
						return
					}
					continue
				}
			}

		case <-interrupt:
			golog.Debug("interrupt")
			if c != nil {
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					golog.Error("id"+strconv.Itoa(id)+"_ Close write error:", err)
					return
				}
			}
			select {
			case <-done:
			case <-time.After(time.Second):
				fmt.Println("wait read time out.Force stop.")
			}
			return
		}
	}
}

func GetRandomToken() string {
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// func newLogFile() *os.File {
// 	filename := "benchTest0416_1.log"
// 	// open an output file, this will append to the today's file if server restarted.
// 	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return f
// }

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
	Cmd     string
	Payload LogData
}

type LogData struct {
	Body string
}

type RespData struct {
	Cmd  string
	Code int
}
