package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dterei/gotsc"
)

func main() {
	for {
		one()
	}
}

func connect() net.Conn {
	for {
		conn, err := net.Dial("tcp", "127.0.0.1:1234")
		if err == nil {
			return conn
		}
		fmt.Println("Error dialing:", err.Error())
		time.Sleep(time.Second)
	}
}

var count = 0

func one() {

	conn := connect()
	defer conn.Close()

	var output *os.File

	r := bufio.NewReader(conn)

	var initialTSCoffset *int64
	var initialLocalTSC *uint64
	var initialRemoteTSC *uint64
	startTime := time.Now()
	for {
		conn.SetDeadline(time.Now().Add(time.Second))
		if _, err := io.WriteString(conn, "\n"); err != nil {
			return
		}
		x, err := r.ReadString('\n')
		if err != nil {
			return
		}
		localTime := time.Now()
		localTSC := gotsc.BenchStart()
		var remoteTime time.Time
		bits := strings.Split(x[0:len(x)-1], " ")
		if err = remoteTime.UnmarshalText([]byte(bits[0])); err != nil {
			log.Fatal("UnmarshalText", err)
		}
		offset, err := strconv.Atoi(bits[1])
		if err != nil {
			log.Fatal("Failed to parse offset", err)
		}
		remoteTSC, err := strconv.ParseUint(bits[2], 10, 64)
		if err != nil {
			log.Fatal("Failed to parse remote TSC", err)
		}
		freq, err := strconv.Atoi(bits[3])
		if err != nil {
			log.Fatal("Failed to parse freq", err)
		}
		state, err := strconv.Atoi(bits[4])
		if err != nil {
			log.Fatal("Failed to parse state", err)
		}

		if initialLocalTSC == nil {
			initialLocalTSC = &localTSC
		}
		if initialRemoteTSC == nil {
			initialRemoteTSC = &remoteTSC
		}
		TSCoffset := int64(remoteTSC - localTSC)
		if initialTSCoffset == nil {
			initialTSCoffset = &TSCoffset
		}

		diff := localTime.Sub(remoteTime)
		// only create the file when we have actual data as vpnkit
		// will accept connections even when the service isn't running
		// yet and we don't want to create lots of empty files.
		if output == nil {
			filename := fmt.Sprintf("drift.%d.dat", count)
			count++
			output, err = os.Create(filename)
			if err != nil {
				log.Fatal("Failed to create output file", err)
			}
			defer output.Close()
			log.Printf("Created %s", filename)
			fmt.Fprintf(output, "# <local time> <delta in microseconds> <TSC delta> <kernel offset> <kernel frequency> <kernel state>\n")
		}
		fmt.Fprintf(output, "%.1f %d %d %d %d %d %d %d\n", localTime.Sub(startTime).Seconds(), diff/1000, localTSC-*initialLocalTSC, remoteTSC-*initialRemoteTSC, TSCoffset-*initialTSCoffset, offset, freq, state)
		time.Sleep(time.Second)
	}
}

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
		conn.Write(bs)
		conn.Write([]byte("\n"))
	}
}
