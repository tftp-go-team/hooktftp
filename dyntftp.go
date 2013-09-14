package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type DataResponse struct {
	data []byte
	bytesRead int
	err error
}

func ReadToChannel(path string, dataChan, returnChan chan *DataResponse) {
	datares := <- returnChan
	file, err := os.Open(path)
	if err != nil {
		datares.err = err
		dataChan <- datares
		return
	}

	for {
		bytesRead, err := file.Read(datares.data)
		datares.bytesRead = bytesRead
		datares.err = err
		dataChan <- datares
		if err != nil {
			if err := file.Close(); err != nil {
				fmt.Println("Failed to close", path, err)
			}
			break
		}
		datares = <- returnChan
	}

}

func SendFile(path string, blocksize int, addr *net.UDPAddr) {
	started := time.Now()

	rrq, err := NewRRQresponse(addr, blocksize)
	if err != nil {
		fmt.Println("Failed to create rrq", err)
		return
	}
	fmt.Println("GET", path, "blocksize", rrq.blocksize)

	if err := rrq.WriteOACK(); err != nil {
		return
	}

	bufsize := 10
	dataChan := make(chan *DataResponse, bufsize)
	returnChan := make(chan *DataResponse, bufsize)

	datapool := make([]byte, blocksize*bufsize)
	for i := 0; i < bufsize; i++ {
		start := blocksize * i
		end := start + blocksize
		returnChan <- &DataResponse{data: datapool[start:end]}
	}

	go ReadToChannel(path, dataChan, returnChan)

	totalBytes := 0

	for {
		datares := <- dataChan
		totalBytes += datares.bytesRead
		if datares.err == io.EOF {
			rrq.Write(datares.data[:datares.bytesRead])
			rrq.End()
			break
		} else if datares.err != nil {
			fmt.Println("Failed to read data from file:", datares.err)
			break
		} else {
			rrq.Write(datares.data[:datares.bytesRead])
		}

		returnChan <- datares
	}

	took := time.Since(started)

	speed := float64(totalBytes) / took.Seconds() / 1024 / 1024

	fmt.Printf("Sent %v bytes in %v %f MB/s\n", totalBytes, took, speed)
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
		fmt.Println("Failed to listen UDP", err)
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
