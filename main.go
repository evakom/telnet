/*
 * HomeWork-10: telnet client
 * Created on 05.11.2019 22:03
 * Copyright (c) 2019 - Eugene Klimov
 */

// Simple telnet client
package main

import (
	"log"
	"time"
)

func main() {

	args := getCmdArgsMap()
	timeout, err := time.ParseDuration(args["timeout"])
	if err != nil {
		log.Fatalln(err)
	}

	client := newClient(args["addr"], timeout)
	if err := client.dial(); err != nil {
		log.Fatalln("Cannot connect:", err)
	}

	abort := client.readFromConn()
	_ = client.writeToConn()
	client.waitOSKill()

	<-abort
	client.cancel()

	time.Sleep(time.Second) // wait deadline for every socket

	if err := client.close(); err != nil {
		log.Fatalln("Error close client:", err)
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
