package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"strconv"
)

var host string

func main () {
	host := "192.168.0.12:8000"
	ln, _ := net.Listen("tcp", host)
	fmt.Println("Listening")
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		defer conn.Close()
		fmt.Println("This is happening")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	bufferIn := bufio.NewReader(conn)
	result, _ := bufferIn.ReadString('\n')
	items := strings.Fields(result)
	// fmt.Println("M1: ", items[1])
	// fmt.Println("M2:", items[2])
	// fmt.Println("Result: ",result)
	fmt.Println("Items: ",items)

	var parsed []float64

	s1, _ := strconv.ParseFloat(items[1], 64)
	s2, _ := strconv.ParseFloat(items[2], 64)
	parsed = append(parsed, s1)
	parsed = append(parsed, s2)
	fmt.Println("Parsed: ",parsed)
	resultFloat := parsed[0]*parsed[1]
	resultString := fmt.Sprintf("%f", resultFloat)
	fmt.Println("Float: ", resultFloat)
	conn.Write([]byte(resultString))
	fmt.Println("String: ", resultString)

}
