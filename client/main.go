package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Following constants should be synced with cilium CI.

const MSG_SIZE = 256
const CLIENT_CONNECTED_MSG = "client successfully connected"
const IO_TIME_OUT = 2 * time.Second

func panicOnErr(ctx string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", ctx, err))
	}
}

func Run(servAddr *net.UDPAddr) {
	var (
		servFirstReply string
		dialConn       *net.UDPConn
		conn           *net.UDPConn
		err            error
	)

	for i := 0; i < 30; i++ {
		// Loop until the connect request goes through once the server is up and running.
		dialConn, err = net.DialUDP("udp", nil, servAddr)
		if err == nil {
			fmt.Printf("%s to %s\n", CLIENT_CONNECTED_MSG, dialConn.RemoteAddr().String())
			dialConn.Close()
			break
		}
		time.Sleep(1 * time.Second)
	}
	panicOnErr("Failed to connect", err)

	dummyAddr, err := net.ResolveUDPAddr("udp", ":0")
	// This is just a dummy call to get net.PacketConn so that we can make unconnected UDP
	// calls to send and receive messages.
	conn, err = net.ListenUDP("udp", dummyAddr)
	panicOnErr("Failed to listen (dummy)", err)
	defer conn.Close()

	request := []byte("hello")
	for {
		reply := make([]byte, MSG_SIZE)

		err = conn.SetWriteDeadline(time.Now().Add(IO_TIME_OUT))
		panicOnErr("SetWriteDeadline", err)
		_, err = conn.WriteToUDP(request, servAddr)
		panicOnErr("Failed to write", err)

		err = conn.SetReadDeadline(time.Now().Add(time.Second))
		panicOnErr("SetReadDeadline", err)
		n, _, err := conn.ReadFromUDP(reply)
		panicOnErr("Failed to read", err)

		fmt.Println("client received: ", string(reply[:n]))

		resStr := string(reply[:n])
		if resStr == "" {
			panic(fmt.Sprintf("Empty response from the server"))

		}

		if servFirstReply == "" {
			servFirstReply = resStr
		} else {
			if resStr != servFirstReply {
				panic(fmt.Sprintf("server reply mismatch new(%s) != old(%s)", resStr, servFirstReply))
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	var (
		err      error
		servAddr *net.UDPAddr
	)
	remote := os.Args[1]

	for i := 0; i < 10; i++ {
		servAddr, err = net.ResolveUDPAddr("udp", remote)
	}
	panicOnErr(fmt.Sprintf("Failed to resolve UDP address[%s]:", remote), err)

	for {
		Run(servAddr)
	}
}
