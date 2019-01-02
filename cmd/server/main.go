package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting server")
	l, err := net.Listen("tcp", "0.0.0.0:1234")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on 0.0.0.0:1234")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

const (
	staNano = 0x2000
)

func handleRequest(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		if _, err := r.ReadString('\n'); err != nil {
			log.Println("Error reading", err)
			return
		}
		bs, err := time.Now().MarshalText()
		if err != nil {
			log.Fatal("Error marshalling time", err)
		}
		var timex syscall.Timex
		state, err := syscall.Adjtimex(&timex)
		if err != nil {
			log.Fatal("Error calling adjtimex", err)
		}
		offset := time.Duration(timex.Offset)
		if timex.Status&staNano == 0 {
			// it's actually microseconds
			offset = time.Duration(timex.Offset * 1000)
		}
		frequency := timex.Freq
		conn.Write(bs)
		conn.Write([]byte(" "))
		conn.Write([]byte(fmt.Sprintf("%d %d %d", int64(offset/time.Microsecond), frequency, state)))
		conn.Write([]byte("\n"))
	}
}
