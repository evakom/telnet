/*
 * HomeWork-10: telnet client
 * Created on 05.11.2019 22:03
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

func main() {

	timeoutArg := flag.String("timeout", "10s", "timeout for connection (duration)")
	fileName := filepath.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Printf("usage: %s [--timeout] <host> <port>\n", fileName)
		fmt.Printf("example1: %s 1.2.3.4 567\n", fileName)
		fmt.Printf("example2: %s --timeout=10s 8.9.10.11 1213\n", fileName)
		flag.PrintDefaults()
	}
	flag.Parse()
	if len(flag.Args()) < 2 {
		flag.Usage()
		os.Exit(2)
	}
	timeout, err := time.ParseDuration(*timeoutArg)
	if err != nil {
		log.Fatalln(err)
	}
	addr := flag.Arg(0) + ":" + flag.Arg(1)

	dialer := &net.Dialer{Timeout: timeout}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}

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
