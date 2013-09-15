package main

import (
	"fmt"
	"net"
	"os"
	"io"
	"time"
	"flag"
	"path/filepath"
	"strings"
)


func SendFile(path string, blocksize int, addr *net.UDPAddr) {
	path = filepath.Join(*root, path)

	started := time.Now()

	rrq, err := NewRRQresponse(addr, blocksize)
	if err != nil {
		fmt.Println("Failed to create rrq", err)
		return
	}

	path, err = filepath.Abs(path)

	if err != nil {
		fmt.Println("Bad path", err)
		rrq.WriteError(UNKNOWN_ERROR, "Invalid file path:" + err.Error())
		return
	}
	fmt.Println("GET", path, "blocksize", rrq.blocksize)

	if !strings.HasPrefix(path, *root) {
		fmt.Println("Path access violation", path)
		rrq.WriteError(ACCESS_VIOLATION, "Path access violation")
		return
	}


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


var root = flag.String("root", "/var/lib/tftpboot/", "Serve files from")
var port = flag.Int("port", 69, "Port to listen")

func main() {
	flag.Parse()
	*root, _ = filepath.Abs(*root)

	fmt.Println("flags", *root, *port)
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
