package main

import (
	"fmt"
	"net"
	"os"
	"io"
	"time"
)


func SendFile(path string, blocksize int, addr *net.UDPAddr) {
	started := time.Now()

	rrq, err := NewRRQresponse(addr, blocksize)
	if err != nil {
		fmt.Println("Failed to create rrq", err)
		return
	}
	fmt.Println("GET", path, "blocksize", rrq.blocksize)

	if err := rrq.WriteOACK(); err != nil {
		fmt.Println("Failed to write OACK", err)
		return
	}


	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Failed to open file", path, err)
		rrq.WriteError(NOT_FOUND, err.Error())
		return
	}

	defer file.Close()

	b := make([]byte, rrq.blocksize)

	totalBytes := 0

	for {
		bytesRead, err := file.Read(b)
		totalBytes += bytesRead

		if err == io.EOF {
			if _, err := rrq.Write(b[:bytesRead]); err != nil {
				fmt.Println("Failed to write last bytes of the file", err)
				return
			}
			rrq.End()
			break
		} else if err != nil {
			fmt.Println("Error while reading", file, err)
			rrq.WriteError(UNKNOWN_ERROR, err.Error())
			return
		}

		if _, err := rrq.Write(b[:bytesRead]); err != nil {
			fmt.Println("Failed to write bytes for", path, err)
			return
		}
	}

	took := time.Since(started)

	speed := float64(totalBytes) / took.Seconds() / 1024 / 1024

	fmt.Printf("Sent %v bytes in %v %f MB/s\n", totalBytes, took, speed)
}


func main() {
	addr, err := net.ResolveUDPAddr("udp", ":1234")
	if err != nil {
		fmt.Println("Failed to resolve address", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Failed to listen UDP" , err)
		return
	}

	data := make([]byte, 50)

	for {
		written, client_addr, err := conn.ReadFrom(data)
		if err != nil {
			fmt.Println("Failed to read data from client:", err)
			continue
		}

		raddr, err := net.ResolveUDPAddr("udp", client_addr.String())
		if err != nil {
			fmt.Println("Failed to resolve client address:", err)
			continue
		}

		request, err := ParseRequest(data[:written])
		if err != nil {
			fmt.Println("Failed to parse request:", err)
			continue
		}

		if request.opcode == RRQ {
			go SendFile(request.path, request.blocksize, raddr)
		} else {
			fmt.Println("Unimplemented opcode:", request.opcode)
		}

	}

}
