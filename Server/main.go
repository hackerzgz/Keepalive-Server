package main

import (
	utils "Keepalive-Server/Utils"
	"log"
	"net"
)

// golang 实现带有心跳检测的 TCP 长链接
// server

// message struct:
// c#d
var (
	Req_REGISTER byte = 1 // 1 --> client register cid
	Res_REGISTER byte = 2 // 2 --> server response

	Req_HEARTBEAT byte = 3 // 3 --> server send heartbeat req
	Res_HEARTBEAT byte = 4 // 4 --> client send heartbeat res

	Req byte = 5 // 5 --> client/server send data
	Res byte = 6 // 6 --> client/server send ack
)

type CS struct {
	Rch chan []byte // Read Chan
	Wch chan []byte // Write Chan
	Dch chan bool   // Done Chan
	u   string      // uid
}

func NewCS(uid string) *CS {
	return &CS{
		Rch: make(chan []byte),
		Wch: make(chan []byte),
		u:   uid,
	}
}

var CMap map[string]*CS

func main() {
	CMap = make(map[string]*CS)

	host, port := utils.GetTCPConf("../TCPConf.json")

	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(host), port, ""})
	if err != nil {
		log.Println("Listen Socket Port", port, "Fail -->", err)
		return
	}

	log.Println("Init Conn Done, Wait for Client...")
	// Do Server
	Server(listen)
}
