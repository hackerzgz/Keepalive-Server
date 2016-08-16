package main

import (
	"fmt"
	"log"
	"net"

	utils "Keepalive-Server/Utils"
)

// golang 实现带有心跳检测的 TCP 长链接
// client

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

var Dch chan bool
var Rch chan []byte
var Wch chan []byte

func main() {
	Dch = make(chan bool)
	Rch = make(chan []byte)
	Wch = make(chan []byte)

	host, port := utils.GetTCPConf("../TCPConf.json")
	tcpAddr := host + ":" + string(port)
	addr, err := net.ResolveTCPAddr("tcp", tcpAddr)
	conn, err := net.DialTCP("tcp", nil, addr)

	//	conn, err := net.Dial("tcp", "127.0.0.1:6666")
	if err != nil {
		log.Println("Conn Server Error -->", err.Error())
		return
	}
	fmt.Println("Server Conn Done.")
	defer conn.Close()

	go Handler(conn)
	select {
	case <-Dch:
		fmt.Println("Conn Close!")
	}
}

func Handler(conn *net.TCPConn) {}
