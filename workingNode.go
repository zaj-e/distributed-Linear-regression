package main

import (
	"bufio"
	"fmt"
	"net"
)

var host string

func main () {
	host := "194.168.0.4:8000"
	ln, _ := net.Listen("tcp", host)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	bufferIn := bufio.NewReader(conn)
	result, _ := bufferIn.ReadString('\n')
	fmt.Println(result)

	newConn, _ := net.Dial("tcp", host)

	defer conn.Close()
	fmt.Println(newConn, result[0]*result[1])
}