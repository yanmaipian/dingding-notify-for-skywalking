package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	//服务器端口号
	serverPort string
	//钉钉告警url
	dingdingUrl string
)

//钉钉告警模板
const template = `{
     "msgtype": "text",
     "text": {
         "content": "%s"
     },
     "at": {
         "isAtAll": false
     }
 }`

func init() {

	flag.StringVar(&serverPort, "p", "9201", "server port")
	flag.StringVar(&dingdingUrl, "u", "", "alarm webhook")
}

func sendMsg(w http.ResponseWriter, r *http.Request) {
	//读取skywalking告警消息
	body, _ := ioutil.ReadAll(r.Body)
	//拼接钉钉所需的消息格式
	msg := strings.NewReader(fmt.Sprintf(template, body))
	//发送钉钉消息
	res, err := http.Post(dingdingUrl, "application/json", msg)

	if err != nil || res.StatusCode != 200 {
		log.Print("send alram msg to dingding fatal.", err, res.StatusCode)
	}

}

func main() {

	flag.Parse()

	if dingdingUrl == "" {
		log.Fatal("dingdingUrl cannot be empty usage -h get help. ")
	}

	http.HandleFunc("/alarm", sendMsg)

	//启动web服务器
	if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), nil); err != nil {
		log.Fatal("server start fatal.", err)
	}

}
