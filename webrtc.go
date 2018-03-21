package main

import (
	"fmt"

	"github.com/tidwall/gjson"

	"golang.org/x/net/websocket"

	"log"

	"net/http"
)

//客户端连接对象
var clientConnectionMap map[string]*websocket.Conn = make(map[string]*websocket.Conn)
var clientCandidateMap map[string]string = make(map[string]string)
var clientOfferMap map[string]string = make(map[string]string)
var clientAnswerMap map[string]string = make(map[string]string)

/**********
	接收到客户端连接
*********/
func mangeServer(ws *websocket.Conn) {

	var err error
	var isConnection = 0
	for {

		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {

			fmt.Println("Can't receive")
			break
		}

		event := gjson.Get(reply, "event").String()
		clientip := gjson.Get(reply, "clientip").String()
		callip := gjson.Get(reply, "callip").String()
		fmt.Println(callip)

		if event == "connection" {
			isConnection = 1
			clientConnectionMap[clientip] = ws
		} else {
			isConnection = 0
		}

		if event == "_ice_candidate" {
			clientCandidateMap[clientip] = reply
		}
		if event == "_offer" {
			//fmt.Println(callip)
			clientOfferMap[clientip] = reply
		}

		if event == "_answer" {
			clientAnswerMap[clientip] = reply
		}

		if isConnection == 0 {
			writeClientData(clientip, reply, event, callip)
		}

	}

}

/*******
	向客户端写入数据
**********/
func writeClientData(ip string, reply string, eventtype string, callip string) {
	var err error
	var clientip string
	var clientWs *websocket.Conn

	for clientip, clientWs = range clientConnectionMap {

		if callip == clientip {

			if err = websocket.Message.Send(clientWs, reply); err != nil {

				fmt.Println("Can't send")

			}
		}
	}

}

/********
	开启websocketserver
*********/
func startWebSocketServer() {
	serverPort := "3000"

	fmt.Println("starting server at port:", serverPort)

	http.Handle("/", websocket.Handler(mangeServer))

	if err := http.ListenAndServe(":"+serverPort, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

/**********
	开启静态文件server
********/
func startStaticServer() {
	h := http.FileServer(http.Dir("html"))
	err := http.ListenAndServe(":9090", h)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {

	go startWebSocketServer()

	startStaticServer()

}
