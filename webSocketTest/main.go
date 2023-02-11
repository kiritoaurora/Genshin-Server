package main

import (
	"webSocketTest/handler"

	"golang.org/x/net/websocket"
)

// type MsgLogin struct {
// 	MsgId    int    `json:"msgId"`
// 	Account  string `json:"account"`
// 	Password string `json:"password"`
// 	UserId   int64  `json:"userId"`
// }

// type MsgPool struct {
// 	MsgId    int `json:"msgId"`
// 	PoolType int `json:"pooltype"`
// }

// type MsgResponsePool struct {
// 	MsgId        int    `json:"msgId"`
// 	ItemId       int    `json:"itemId"`
// 	ItemName     string `json:"itemName"`
// 	Stuff        string `json:"stuff"`
// 	StuffNum     int64  `json:"stuffNum"`
// 	StuffItem    string `json:"stuffItem"`
// 	StuffItemNum int64  `json:"stuffItemNum"`
// }

func main() {
	ws, err := websocket.Dial("ws://127.0.0.1:8888/", "", "http://127.0.0.1:8888/")
	if err != nil {
		return
	}
	defer ws.Close()

	go handler.ListenMsg(ws)

	handler.Login(ws)
	handler.Run(ws)
}
