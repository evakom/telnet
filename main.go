/*
 * HomeWork-10: telnet client
 * Created on 05.11.2019 22:03
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"bufio"
	"context"
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

	ctx, cancel := context.WithCancel(context.Background())

	dialer := &net.Dialer{Timeout: timeout}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		log.Fatalln("Cannot connect:", err)
	}
	fmt.Println("Connected to:", addr)
	fmt.Println("Press 'Ctrl+D' for exit")

	go readRoutine(ctx, conn)
	go writeRoutine(ctx, conn)

	time.Sleep(10 * time.Second)

	fmt.Println("Cancelling all operations... ")
	cancel()
	time.Sleep(1 * time.Second) // 0.5 second for every socket goroutine
	fmt.Println("...canceled all operations")

	fmt.Println("Closing connection... ")
	if err := conn.Close(); err != nil {
		log.Fatalln("Error close connection:", err)
	}
	fmt.Println("...closed connection")
	fmt.Println("Exited.")
}

func readRoutine(ctx context.Context, conn net.Conn) {
	reply := make([]byte, 1)
OUTER:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Cancel happened, exiting from reading...")
			break OUTER
		default:
			// set deadline for read socket - need 'select loop' continue
			err := conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
			if err != nil {
				log.Println(err)
			}
			n, err := conn.Read(reply)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && !netErr.Timeout() {
					log.Println(err)
				}
				//if err == io.EOF {
				//	fmt.Println("dddddddddddddddddddddddddd", err)
				//}
			}
			if n == 0 {
				break
			}
			fmt.Print(string(reply))
		}
	}
	fmt.Println("...exited from reading")
}

func writeRoutine(ctx context.Context, conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				break OUTER
			}
			str := scanner.Text()
			log.Printf("To server %v\n", str)

			conn.Write([]byte(fmt.Sprintf("%s\n", str)))
		}

	}
	log.Printf("Finished writeRoutine")
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
