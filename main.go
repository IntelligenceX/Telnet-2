/*
Use: Telnet-2 [Host] [Port] [Optional Timeout]

Timeout (in seconds) is optional, default is 10 seconds. Host and Port are required.
By default, it uses the local proxy "127.0.0.1:9050" where Tor is expected to run. If Tor is not needed, change the proxy field and recompile.
*/

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// proxy to connect to. Empty to disable.
const proxyAddress = "127.0.0.1:9050"

func main() {
	// parse the flags: Telnet-2 [Host] [Port] [Optional Timeout]
	if len(os.Args) < 3 {
		fmt.Println("Use: Telnet-2 [Host] [Port] [Optional Timeout]")
		return
	}

	host := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid port")
		return
	}

	timeout := time.Second * 10
	if len(os.Args) > 3 {
		if timeout2, err := strconv.Atoi(os.Args[3]); err != nil {
			timeout = time.Duration(timeout2)
		}
	}

	// prepare the client
	telnetClient := NewTelnetClient(host, port, timeout, proxyAddress)
	telnetClient.ProcessData(os.Stdin, os.Stdout)
}
