package main

import (
	"fmt"
	"net"
)


func SendFile(path string, blocksize int, addr *net.UDPAddr) {

	rrq, err := NewRRQresponse(addr, blocksize)
	if err != nil {
		fmt.Println("Failed to create rrq", err)
		return
	}

	rrq.WriteOACK()
	rrq.Write([]byte("Hello World Foo Bar jea end"))

	rrq.End()

}


func main() {
	fmt.Println("hello")
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
		_, client_addr, err := conn.ReadFrom(data)
		if err != nil {
			fmt.Println("Failed to read data from client:", err)
			continue
		}

		raddr, err := net.ResolveUDPAddr("udp", client_addr.String())
		if err != nil {
			fmt.Println("Failed to resolve client address:", err)
			continue
		}

		request, err := ParseRequest(data)
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
