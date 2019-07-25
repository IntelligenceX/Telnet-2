// Fork from https://github.com/mtojek/go-telnet but heavily modified to simplify and to optionally use a proxy.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"golang.org/x/net/proxy"
)

const defaultBufferSize = 4096

// TelnetClient represents a TCP client which is responsible for writing input data and printing response.
type TelnetClient struct {
	destination     *net.TCPAddr
	responseTimeout time.Duration
	proxy           string // Proxy to use (in the form of "IP:Port"), empty to disable
}

// NewTelnetClient method creates new instance of TCP client.
// Proxy is in the form "IP:Port" or empty to disable.
func NewTelnetClient(Host string, Port int, Timeout time.Duration, Proxy string) *TelnetClient {
	tcpAddr := net.JoinHostPort(Host, strconv.Itoa(Port))
	resolved := resolveTCPAddr(tcpAddr)

	return &TelnetClient{
		destination:     resolved,
		responseTimeout: Timeout,
		proxy:           Proxy,
	}
}

func resolveTCPAddr(addr string) *net.TCPAddr {
	resolved, error := net.ResolveTCPAddr("tcp", addr)
	if nil != error {
		log.Fatalf("Error occured while resolving TCP address \"%v\": %v\n", addr, error)
	}

	return resolved
}

// ProcessData method processes data: reads from input and writes to output.
func (t *TelnetClient) ProcessData(inputData io.Reader, outputData io.Writer) {
	var connection net.Conn
	var err error

	// optionally use the proxy, if set
	if t.proxy != "" {
		var dialSocksProxy proxy.Dialer
		dialSocksProxy, err = proxy.SOCKS5("tcp", t.proxy, nil, proxy.Direct)
		if err != nil {
			return
		}
		connection, err = dialSocksProxy.Dial("tcp", t.destination.String())
	} else {
		connection, err = net.DialTCP("tcp", nil, t.destination)
	}

	if err != nil {
		log.Fatalf("Error occured while connecting to address \"%v\": %v\n", t.destination.String(), err)
	}

	defer connection.Close()

	requestDataChannel := make(chan []byte)
	doneChannel := make(chan bool)
	responseDataChannel := make(chan []byte)

	go t.readInputData(inputData, requestDataChannel, doneChannel)
	go t.readServerData(connection, responseDataChannel)

	var afterEOFResponseTicker = new(time.Ticker)
	var afterEOFMode bool
	var somethingRead bool

	for {
		select {
		case request := <-requestDataChannel:
			if _, error := connection.Write(request); nil != error {
				log.Fatalf("Error occured while writing to TCP socket: %v\n", error)
			}
		case <-doneChannel:
			afterEOFMode = true
			afterEOFResponseTicker = time.NewTicker(t.responseTimeout)
		case response := <-responseDataChannel:
			outputData.Write([]byte(fmt.Sprintf("%v", string(response))))
			somethingRead = true

			if afterEOFMode {
				afterEOFResponseTicker.Stop()
				afterEOFResponseTicker = time.NewTicker(t.responseTimeout)
			}
		case <-afterEOFResponseTicker.C:
			if !somethingRead {
				log.Println("Nothing read. Maybe connection timeout.")
			}
			return
		}
	}
}

func (t *TelnetClient) readInputData(inputData io.Reader, toSent chan<- []byte, doneChannel chan<- bool) {
	buffer := make([]byte, defaultBufferSize)
	var error error
	var n int

	reader := bufio.NewReader(inputData)

	for nil == error {
		n, error = reader.Read(buffer)
		toSent <- buffer[:n]
	}

	t.assertEOF(error)
	doneChannel <- true
}

func (t *TelnetClient) readServerData(connection net.Conn, received chan<- []byte) {
	buffer := make([]byte, defaultBufferSize)
	var error error
	var n int

	for nil == error {
		n, error = connection.Read(buffer)
		received <- buffer[:n]
	}

	t.assertEOF(error)
}

func (t *TelnetClient) assertEOF(err error) {
	if "EOF" != err.Error() {
		log.Fatalf("Error occured while operating on TCP socket: %v\n", err)
	}
}
