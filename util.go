/*
 * HomeWork-10: telnet client
 * Created on 08.11.2019 19:44
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type args map[string]string

func getCmdArgsMap() args {
	args := make(args)

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

	args["addr"] = flag.Arg(0) + ":" + flag.Arg(1)
	args["timeout"] = *timeoutArg

	return args
}
