/*
 * HomeWork-10: telnet client
 * Created on 08.11.2019 22:24
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

const (
	SERVERLISTEN      = "localhost:12345"
	SERVERSTOPMESSAGE = "stopServer\n"
	SERVERWAITSTART   = 100 * time.Millisecond
)

func init() {
	//f, err := os.OpenFile("client_test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//log.SetOutput(f)
	// or
	//log.SetOutput(ioutil.Discard)
}

func TestDialAndClose(t *testing.T) {
	go startServer()
	time.Sleep(SERVERWAITSTART)

	client := newClient(SERVERLISTEN, 10*time.Nanosecond)
	if err := client.dial(); err == nil {
		t.Fatal("Client successfully connected with small timeout 10ns but expected i/o error")
	}

	client = newClient(SERVERLISTEN, 10*time.Second)
	if err := client.dial(); err != nil {
		t.Fatalf("Expected successfully connected to server but got error: %s", err)
	}

	time.Sleep(1 * time.Second)

	if err := client.close(); err != nil {
		t.Fatalf("Expected successfully closed connection to server but got error: %s", err)
	}

	stopServer()
}

func startServer() {

	ln, err := net.Listen("tcp", SERVERLISTEN)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Test server started...")

	conn, err := ln.Accept()
	if err != nil {
		log.Fatalln(err)
	}

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil && err != io.EOF {
			log.Fatalln(err)
		}
		if message == SERVERSTOPMESSAGE || err == io.EOF {
			break
		}
		answer := strings.ToUpper(message)
		if _, err = conn.Write([]byte(answer)); err != nil {
			log.Fatalln(err)
		}
	}
	if err := conn.Close(); err != nil {
		log.Fatalln(err)
	}

	log.Println("...test server stopped.")
}

func stopServer() {
	dialer := &net.Dialer{}
	conn, err := dialer.Dial("tcp", SERVERLISTEN)
	if err != nil {
		log.Fatalln(err)
	}
	if _, err = conn.Write([]byte(SERVERSTOPMESSAGE)); err != nil {
		log.Fatalln(err)
	}
}
