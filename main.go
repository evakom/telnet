/*
 * HomeWork-10: telnet client
 * Created on 05.11.2019 22:03
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:1230")
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(conn, text+"\n")
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}

//func main() {
//
//	fmt.Println("Launching server...")
//
//	// listen on all interfaces
//	ln, _ := net.Listen("tcp", ":8081")
//
//	// accept connection on port
//	conn, _ := ln.Accept()
//
//	// run loop forever (or until ctrl-c)
//	for {
//		// will listen for message to process ending in newline (\n)
//		message, _ := bufio.NewReader(conn).ReadString('\n')
//		// output message received
//		fmt.Print("Message Received:", string(message))
//		// sample process for string received
//		newmessage := strings.ToUpper(message)
//		// send new string back to client
//		conn.Write([]byte(newmessage + "\n"))
//	}
//}
