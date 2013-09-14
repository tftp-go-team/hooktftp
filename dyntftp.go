package main

import (
	"fmt"
	"net"
	"os"
	"io"
)


func SendFile(path string, blocksize int, addr *net.UDPAddr) {

	rrq, err := NewRRQresponse(addr, blocksize)
	if err != nil {
		fmt.Println("Failed to create rrq", err)
		return
	}
	fmt.Println("GET", path, "blocksize", rrq.blocksize)

	if err := rrq.WriteOACK(); err != nil {
		return
	}


	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Failed to open file", path, err)
		// TODO: write error package
		return
	}

	b := make([]byte, rrq.blocksize)


	for {
		bytesRead, err := file.Read(b)

		if err == io.EOF {
			rrq.Write(b[:bytesRead])
			rrq.End()
			break
		} else if err != nil {
			fmt.Println("Error while reading", file, err)
			// TODO: write error package
			return
		}

		rrq.Write(b[:bytesRead])
	}

	file.Close()
	fmt.Println("END-GET", path)
	fmt.Println()
	fmt.Println()

}


func main() {
	fmt.Println("hello2")
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
