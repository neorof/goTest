package main

import (
	"fmt"
	"math/rand"
	"time"
)

const str = "0123456789abcdefghijklmnopqrstuvwxyz"

const logDa1 = "{\"body\" : \"wxapp_type=xcx&platform=wxapp&app=pdd&page_name=index&page_id=10002_1551430328032_8V4bP3wwBL&page_sn=10002&"
const logDa2 = "page_url=pages%2Findex%2Findex&page=index&op=epv&sub_op=leave&page_duration=66345&enter_time=1551430328084&user_id=6677178664&log_id=15514303944293t0jBYCVdLCvQHqd&token=LKIG7BPWHL7TPPEORCOOEG6SCG3JPCTLG74CIXJPM4LII5X7NJ3Q101287d&app_version=v3.1.6.2&time=1551430394299&xcx_x_scene=1089&pdd_user_type=0&xcx_trace_id=bQH0A18KIhqkBqZSvBjTDnscsLTZPjBB&session_id=8NLzGS5Vuq9eNmCLvj82ZlhuvdegJoC7&openId=ogJUI0VzqwbaKmprbUy4t7WQyfw0&greyscale_guid=cf1471a9a152de6fe27e4d23423c7eda&withsocket=1\"}"

func main() {
	fmt.Println(logDa1 + GetRandomToken() + "=" + GetRandomToken() + "&" + logDa2)
}

func GetRandomToken() string {
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
