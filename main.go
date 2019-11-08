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

	abort := client.readFromWriteToConn()
	client.waitOSKill()

	<-abort
	client.cancel()

	time.Sleep(time.Second) // wait deadline for every socket

	if err := client.close(); err != nil {
		log.Fatalln("Error close client:", err)
	}
}
