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
	SERVERLISTEN       = "localhost:12346"
	SERVERWAITSTART    = 100 * time.Millisecond
	SERVERWAITSTOP     = 200 * time.Millisecond
	CLIENTTIMEOUT      = 10 * time.Second
	TESTMESSAGE        = "Привет, Otus!\n"
	CLIENTREADTIMEOUT  = 250 * time.Millisecond
	CLIENTCLOSETIMEOUT = 500 * time.Millisecond
)

type server struct {
	conn net.Conn
}

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
	serverOk := newServer()
	go serverOk.startServer()
	time.Sleep(SERVERWAITSTART)

	client := newClient(SERVERLISTEN, 10*time.Nanosecond)
	if err := client.dial(); err == nil {
		t.Error("Client successfully connected with small timeout 10ns but expected i/o error")
	}

	client = newClient(SERVERLISTEN, CLIENTTIMEOUT)
	if err := client.dial(); err != nil {
		t.Errorf("Expected successfully connected to server but got error: %s", err)
	}

	time.Sleep(1 * time.Second)

	if err := client.close(); err != nil {
		t.Errorf("Expected successfully closed connection to server but got error: %s", err)
	}

	serverOk.stopServer()
}

func TestMessageFromServer(t *testing.T) {
	serverOk := newServer()
	go serverOk.startServer()
	time.Sleep(SERVERWAITSTART)

	client := newClient(SERVERLISTEN, CLIENTTIMEOUT)
	if err := client.dial(); err != nil {
		t.Errorf("Expected successfully connected to server but got error: %s", err)
	}

	client.readFromConn()
	serverOk.writeString(TESTMESSAGE)
	time.Sleep(CLIENTREADTIMEOUT)

	if client.lastMessage != TESTMESSAGE {
		t.Errorf("Test message from server no equal client received:\n"+
			"from server - %s, expected - %s", client.lastMessage, TESTMESSAGE)
	}

	_ = client.cancelReadWriteClose()
	time.Sleep(CLIENTCLOSETIMEOUT)

	serverOk.stopServer()
}

func TestMessageToServerCtrlD(t *testing.T) {
	serverOk := newServer()
	go serverOk.startServer()
	time.Sleep(SERVERWAITSTART)

	client := newClient(SERVERLISTEN, CLIENTTIMEOUT)
	if err := client.dial(); err != nil {
		t.Errorf("Expected successfully connected to server but got error: %s", err)
	}

	_ = client.readFromWriteToConn()
	client.stdinChan <- TESTMESSAGE
	time.Sleep(CLIENTREADTIMEOUT)

	actual := client.lastMessage
	expected := strings.ToUpper(TESTMESSAGE)
	if actual != expected {
		t.Errorf("Test message from server expected in upper case:\n"+
			"from server - %s, actual - %s", expected, actual)
	}

	_ = client.cancelReadWriteClose()
	serverOk.stopServer()
}

func TestDisconnectFromServer(t *testing.T) {
	serverOk := newServer()
	go serverOk.startServer()
	time.Sleep(SERVERWAITSTART)

	client := newClient(SERVERLISTEN, CLIENTTIMEOUT)
	if err := client.dial(); err != nil {
		t.Errorf("Expected successfully connected to server but got error: %s", err)
	}

	abort := client.readFromConn()

	serverOk.stopServer()
	time.Sleep(SERVERWAITSTOP * 2)

	<-abort
	_ = client.cancelReadWriteClose()

	if err := client.close(); err == nil {
		t.Error("Expected client err for closing connection but no error returned")
	}
}

func newServer() *server {
	return &server{}
}

func (srv *server) startServer() {

	ln, err := net.Listen("tcp", SERVERLISTEN)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Test server started...")

	srv.conn, err = ln.Accept()
	if err != nil {
		log.Fatalln(err)
	}

	for {
		message, err := bufio.NewReader(srv.conn).ReadString('\n')
		if err != nil && err != io.EOF {
			//log.Println(err)
			// return if closed conn, no need error
			break
		}
		answer := strings.ToUpper(message)
		_, _ = srv.conn.Write([]byte(answer))
		//if err != nil {
		//log.Println(err)
		// return if closed conn, no need error
		//}
	}
}

func (srv *server) stopServer() {
	time.Sleep(SERVERWAITSTOP)
	if err := srv.conn.Close(); err != nil {
		log.Println(err)
	}
	log.Println("...test server stopped.")
}

func (srv *server) writeString(s string) {
	_, _ = srv.conn.Write([]byte(s))
}
