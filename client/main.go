package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Following constants should be synced with cilium CI.

const MSG_SIZE = 256
const IO_TIME_OUT = 5 * time.Second

func panicOnErr(ctx string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", ctx, err))
	}
}

func Run(servAddr *net.UDPAddr) {
	var (
		servFirstReply string
		conn           *net.UDPConn
		err            error
	)

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

		err = conn.SetReadDeadline(time.Now().Add(IO_TIME_OUT))
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
		if servAddr, err = net.ResolveUDPAddr("udp", remote); err != nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	panicOnErr(fmt.Sprintf("Failed to resolve UDP address[%s]:", remote), err)

	for {
		Run(servAddr)
	}
}
