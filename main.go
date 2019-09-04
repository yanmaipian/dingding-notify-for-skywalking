package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	//return string(body)
	data := string(body)
	projstr := gjson.Get(data, "#.name").String()
	mesgstr := gjson.Get(data, "#.alarmMessage").String()
	fn := func(c rune) bool {
		return strings.ContainsRune("[]\"", c)
	}
	mesgstr = strings.TrimFunc(mesgstr, fn)
	projstr = strings.TrimFunc(projstr, fn)
	str := projstr + " " + mesgstr
	fmt.Println(str)
	//拼接钉钉所需的消息格式
	msg := strings.NewReader(fmt.Sprintf(template, str))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: tr}
	req, err := http.NewRequest("POST", dingdingUrl, msg)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Print(err)
		os.Exit(0)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
		os.Exit(0)
	}
	defer res.Body.Close()
	resbody, err := ioutil.ReadAll(res.Body)
	response := string(resbody)
	if err != nil {
		log.Print(err)
		os.Exit(0)
	}
	fmt.Println(response)
}

func main() {

	flag.Parse()
	log.Println("DINGURL", dingdingUrl)
	if dingdingUrl == "" {
		log.Fatal("dingdingUrl cannot be empty usage -h get help. ")
	}

	http.HandleFunc("/alarm", sendMsg)

	//启动web服务器
	if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), nil); err != nil {
		log.Fatal("server start fatal.", err)
	}

}
