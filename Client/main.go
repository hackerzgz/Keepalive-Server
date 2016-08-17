package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

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
	tcpAddr := host + ":" + strconv.Itoa(port)
	fmt.Println("tcpAddr --> ", tcpAddr)
	addr, err := net.ResolveTCPAddr("tcp", tcpAddr)
	conn, err := net.DialTCP("tcp", nil, addr)

	//	conn, err := net.Dial("tcp", "127.0.0.1:6666")
	if err != nil {
		log.Println("Conn Server Error -->", err.Error())
		return
	}
	log.Println("Server Conn Done.")
	defer conn.Close()

	go Handler(conn)
	select {
	case <-Dch:
		log.Println("Conn Close!")
	}
}

// Handler Handle Conn Data.
func Handler(conn *net.TCPConn) {
	// 直到 register ok
	data := make([]byte, 128)

	for {
		conn.Write([]byte{Req_REGISTER, '#', '2'})
		conn.Read(data)

		if data[0] == Res_REGISTER {
			break
		}
	}

	go RHandler(conn)
	go WHandler(conn)
	go Work()
}

// RHandler 心跳包，回复 ack.
func RHandler(conn *net.TCPConn) {
	for {
		data := make([]byte, 128)

		i, _ := conn.Read(data)
		if i == 0 {
			Dch <- true
			return
		}
		if data[0] == Req_HEARTBEAT {
			log.Println("Recv HeartBeat Pack")
			conn.Write([]byte{Res_REGISTER, '#', 'h'})
			log.Println("Send HeartBeat Pack Ack")
		} else if data[0] == Req {
			log.Println("Recv Data Pack")
			fmt.Printf("%v\n", string(data[2:]))
			Rch <- data[2:]
			conn.Write([]byte{Res, '#'})
		}
	}

}

// WHandler 写入数据
func WHandler(conn *net.TCPConn) {
	for {
		select {
		case msg := <-Wch:
			fmt.Println((msg[0]))
			fmt.Println("Send Data After: " + string(msg[1:]))
			conn.Write(msg)
		}
	}
}

// Work 读取数据到写入通道
func Work() {
	for {
		select {
		case msg := <-Rch:
			fmt.Println("Work Recv " + string(msg))
			Wch <- []byte{Req, '#', 'x', 'x', 'x', 'x', 'x'}
		}
	}
}
