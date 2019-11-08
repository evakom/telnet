/*
 * HomeWork-10: telnet client
 * Created on 08.11.2019 19:17
 * Copyright (c) 2019 - Eugene Klimov
 */

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

// DEADLINETIME constant need for waiting user for some time
// and get some work if user "dead"
const DEADLINETIME = time.Millisecond * 500

type client struct {
	serverAddr string
	timeout    time.Duration
	conn       net.Conn
	ctx        context.Context
	cancel     context.CancelFunc
	abort      chan bool
	stdin      chan string
}

func newClient(serverAddr string, timeout time.Duration) client {
	c := client{
		serverAddr: serverAddr,
		timeout:    timeout,
		abort:      make(chan bool),
		stdin:      make(chan string),
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return c
}

func (c *client) dial() error {
	var err error
	dialer := &net.Dialer{Timeout: c.timeout}
	c.conn, err = dialer.Dial("tcp", c.serverAddr)
	if err == nil {
		log.Printf("Connected to: %s", c.serverAddr)
		fmt.Println("Press 'Ctrl+D or Ctrl+C' for exit")
	}
	return err
}

func (c *client) waitOSKill() {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		sig := <-ch
		fmt.Println("\nGot signal:", sig)
		c.abort <- true
	}()
}

func (c *client) close() error {
	log.Print("Closing connection... ")
	if err := c.conn.Close(); err != nil {
		return err
	}
	log.Print("...closed connection")
	log.Print("Exited.")
	return nil
}

func (c *client) readFromConn() chan bool {
	go c.readRoutine()
	return c.abort
}

func (c *client) writeToConn() chan bool {
	go c.writeRoutine()
	return c.abort
}

func (c *client) readFromWriteToConn() chan bool {
	go c.readRoutine()
	go c.writeRoutine()
	return c.abort
}

func (c *client) readRoutine() {
	reply := make([]byte, 1)
OUTER:
	for {
		select {
		case <-c.ctx.Done():
			log.Print("Exiting from reading...")
			break OUTER
		default:
			// set deadline for read socket - need for 'select loop' continue
			if err := c.conn.SetReadDeadline(time.Now().Add(DEADLINETIME)); err != nil {
				log.Println(err)
			}
			n, err := c.conn.Read(reply)
			if err != nil {
				if err == io.EOF {
					log.Println("Remote host aborted connection, exiting from reading...")
					c.abort <- true
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
	log.Println("...exited from reading")
}

func (c *client) writeRoutine() {
	go func(stdin chan<- string) {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					log.Print("Ctrl+D detected, aborting...")
					c.abort <- true
					return
				}
				log.Println(err)
			}
			stdin <- s
		}
	}(c.stdin)

OUTER:
	for {
		select {
		case <-c.ctx.Done():
			log.Print("Exiting from writing...")
			break OUTER
		default:

		STDIN:
			for {
				select {
				case stdin, ok := <-c.stdin:
					if !ok {
						break STDIN
					}
					if _, err := c.conn.Write([]byte(stdin)); err != nil {
						log.Println(err)
					}
					// wait deadline for input
				case <-time.After(DEADLINETIME):
					break STDIN
				}
			}
		}
	}
	log.Println("...exited from writing")
}
