/*
 * HomeWork-10: telnet client
 * Created on 08.11.2019 19:17
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"context"
	"net"
	"time"
)

type client struct {
	serverAddr string
	timeout    time.Duration
	conn       net.Conn
	ctx        context.Context
	cancel     context.CancelFunc
}

func newClient(serverAddr string, timeout time.Duration) client {
	c := client{
		serverAddr: serverAddr,
		timeout:    timeout,
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return c
}

func (c *client) dial() error {
	var err error
	dialer := &net.Dialer{Timeout: c.timeout}
	c.conn, err = dialer.Dial("tcp", c.serverAddr)
	return err
}
