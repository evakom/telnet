/*
 * HomeWork-10: telnet client
 * Created on 05.11.2019 22:03
 * Copyright (c) 2019 - Eugene Klimov
 */

// Simple telnet client
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const DEADLINETIME = time.Millisecond * 500

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
	fmt.Println("Connected to:", args["addr"])
	fmt.Println("Press 'Ctrl+D or Ctrl+C' for exit")

	abort := make(chan bool)
	stdin := make(chan string)

	go readRoutine(client.ctx, client.conn, abort)
	go writeRoutine(client.ctx, client.conn, stdin, abort)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		sig := <-c
		fmt.Println("Got signal:", sig)
		abort <- true
	}()

	<-abort
	client.cancel()

	time.Sleep(DEADLINETIME * 2) // wait DEADLINETIME for every socket goroutine

	fmt.Println("Closing connection... ")
	if err := client.conn.Close(); err != nil {
		log.Fatalln("Error close connection:", err)
	}
	fmt.Println("...closed connection")
	fmt.Println("Exited.")
}

func readRoutine(ctx context.Context, conn net.Conn, abort chan bool) {
	reply := make([]byte, 1)
OUTER:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Exiting from reading...")
			break OUTER
		default:
			// set deadline for read socket - need for 'select loop' continue
			if err := conn.SetReadDeadline(time.Now().Add(DEADLINETIME)); err != nil {
				log.Println(err)
			}
			n, err := conn.Read(reply)
			if err != nil {
				if err == io.EOF {
					fmt.Println("Remote host aborted connection, exiting from reading...")
					abort <- true
					break OUTER
				}
				if netErr, ok := err.(net.Error); ok && !netErr.Timeout() {
					log.Println(err)
				}
			}
			if n == 0 {
				break
			}
			fmt.Print(string(reply))
		}
	}
	fmt.Println("...exited from reading")
}

func writeRoutine(ctx context.Context, conn net.Conn, stdin chan string, abort chan bool) {
	go func(stdin chan<- string) {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Println("Ctrl+D detected, aborting...")
					abort <- true
					return
				}
				log.Println(err)
			}
			stdin <- s
		}
	}(stdin)

OUTER:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Exiting from writing...")
			break OUTER
		default:

		STDIN:
			for {
				select {
				case stdin, ok := <-stdin:
					if !ok {
						break STDIN
					}
					if _, err := conn.Write([]byte(stdin)); err != nil {
						log.Println(err)
					}
					// wait deadline for input
				case <-time.After(DEADLINETIME):
					break STDIN
				}
			}
		}
	}
	fmt.Println("...exited from writing")
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
