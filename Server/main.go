package main

import (
	utils "Keepalive-Server/Utils"
	"fmt"
	"log"
	"net"
	"time"
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

// CS Client/Server
type CS struct {
	Rch chan []byte // Read Chan
	Wch chan []byte // Write Chan
	Dch chan bool   // Done Chan
	u   string      // uid
}

// NewCS return struct CS
// @param uid string
// @return *CS
func NewCS(uid string) *CS {
	return &CS{
		Rch: make(chan []byte),
		Wch: make(chan []byte),
		u:   uid,
	}
}

// CMap Save Message Queue
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
	go PutGRT()
	Server(listen)
}

// PutGRT Put Message to Client Every 15 Seconds.
func PutGRT() {
	for {
		time.Sleep(15 * time.Second)

		for k, v := range CMap {
			fmt.Println("Push Message to User: ", k)
			v.Wch <- []byte{Req, '#', 'p', 'u', 's', 'h', '!'}
		}
	}
}

// Server Listen to Req From Client.
func Server(listen *net.TCPListener) {
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("Accept Client Conn Abnormal --> ", err.Error())
			continue
		}

		fmt.Println("Accept Client From -->", conn.RemoteAddr().String())
		// Handle Goroutine
		go Handle(conn)
	}
}

// Handle Handle Conn Data.
func Handle(conn net.Conn) {
	defer conn.Close()
	data := make([]byte, 128)
	var uid string
	var C *CS

	for {
		conn.Read(data)
		log.Println("Data From Client -->", string(data))

		if data[0] == Req_REGISTER { // register
			conn.Write([]byte{Res_REGISTER, '#', 'o', 'k'})
			uid = string(data[2:])
			C = NewCS(uid)
			CMap[uid] = C
			//			fmt.Println("register client")
			//			fmt.Println(uid)
			break
		} else {
			conn.Write([]byte{Res_REGISTER, '#', 'e', 'r'})
		}
	}

	// Write Handler
	go WHandler(conn, C)

	// Read Handler
	go RHandler(conn, C)

	// Worker
	go Worker(C)

	select {
	case <-C.Dch:
		log.Println("Close Handler Goroutine.")
	}
}

// WHandler 写数据
// 定时 20s 监测超时    conn die ==> goroutine die
func WHandler(conn net.Conn, C *CS) {
	// 读取业务 Work 写入 Wch 的数据
	ticker := time.NewTicker(20 * time.Second)
	for {
		select {
		case d := <-C.Wch:
			conn.Write(d)
		case <-ticker.C:
			// 判断是否还有该用户的数据在队列中
			if _, ok := CMap[C.u]; !ok {
				fmt.Println("Conn Done, Close WHandler")
				return
			}
		}
	}
}

// RHandler 读取客户端数据 ＋ 心跳检测
func RHandler(conn net.Conn, C *CS) {
	//  心跳 ack
	//  将业务数据写入 Wch

	for {
		data := make([]byte, 128)
		// Set Read Timeout.
		err := conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			log.Println("Conn Set Read Timeout Error --> ", err.Error())
		}

		if _, derr := conn.Read(data); derr == nil {
			// 确认是否来自 Client 的 Message
			//      数据消息
			fmt.Println(data)
			if data[0] == Res {
				fmt.Println("Recv Client Data Ack")
			} else if data[0] == Req {
				fmt.Println("Recv Client Data")
				fmt.Println(data)
				conn.Write([]byte{Res, '#'})
				// C.Rch <- data
			}

			continue
		}

		conn.Write([]byte{Req_HEARTBEAT, '#'})
		log.Println("Send HeartBeat Packet")
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		// 检测 Client 是否断开连接
		if _, herr := conn.Read(data); herr == nil {
			log.Println("Resv HeartBeat Packet Ack")
		} else {
			delete(CMap, C.u)
			log.Println("Lose a Client, Delete User!")
			return
		}
	}
}

// Worker 整套工作流程
func Worker(C *CS) {
	time.Sleep(5 * time.Second)
	C.Wch <- []byte{Req, '#', 'h', 'e', 'l', 'l', 'o'}

	time.Sleep(15 * time.Second)
	C.Wch <- []byte{Req, '#', 'h', 'e', 'l', 'l', 'o'}
	// 从 Rch 读信息
	/*	ticker := time.NewTicker(20 * time.Second)
		for {
			select {
			case d := <-C.Rch:
				C.Wch <- d
			case <-ticker.C:
				if _, ok := CMap[C.u]; !ok {
					return
				}
			}

		}
	*/ // 往 Wch 写信息
}
