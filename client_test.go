/*
 * HomeWork-10: telnet client
 * Created on 08.11.2019 22:24
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

const (
	SERVERLISTEN      = "localhost:12345"
	SERVERSTOPMESSAGE = "stopServer\n"
)

func init() {
	//f, err := os.OpenFile("client_test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//log.SetOutput(f)
	log.SetOutput(ioutil.Discard)
}

func TestDial(t *testing.T) {
	go startServer()
	time.Sleep(100 * time.Millisecond)

	client := newClient(SERVERLISTEN, 10*time.Second)
	if err := client.dial(); err != nil {
		log.Fatalln("Cannot connect:", err)
	}

	time.Sleep(10 * time.Second)
	//stopServer()
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
		if err != nil {
			log.Fatalln(err)
		}
		if message == SERVERSTOPMESSAGE {
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

//func stopServer() {
//	dialer := &net.Dialer{}
//	conn, err := dialer.Dial("tcp", SERVERLISTEN)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	if _, err = conn.Write([]byte(SERVERSTOPMESSAGE)); err != nil {
//		log.Fatalln(err)
//	}
//}
