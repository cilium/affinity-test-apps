package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// Following constants should be synced with cilium CI.

const MSG_SIZE = 256
const IO_TIME_OUT = 2 * time.Second

func panicOnErr(ctx string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", ctx, err))
	}
}

func Serve(conn *net.UDPConn, hostname string) {
	response := []byte(hostname)

	for {
		buf := make([]byte, MSG_SIZE)

		_, addr, err := conn.ReadFrom(buf)
		fmt.Printf("received request from %s\n", addr)

		err = conn.SetWriteDeadline(time.Now().Add(IO_TIME_OUT))
		panicOnErr("SetWriteDeadline", err)

		fmt.Printf("Server sent %s \n", string(response))
		_, err = conn.WriteTo(response, addr)
		panicOnErr("Failed to write", err)
	}
}

func main() {
	port := os.Args[1]
	addr := ":" + port

	hostname, err := os.Hostname()
	panicOnErr("hostname error", err)
	log.Printf("hostname %s\n", hostname)

	servAddr, err := net.ResolveUDPAddr("udp", addr)
	panicOnErr(fmt.Sprintf("Failed to resolve UDP address[%s]:", addr), err)
	conn, err := net.ListenUDP("udp", servAddr)
	panicOnErr(fmt.Sprintf("Failed to listen %s", addr), err)
	fmt.Printf("UDP server listening on %s\n", conn.LocalAddr().String())
	defer conn.Close()

	for {
		Serve(conn, hostname)
	}
}
