package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	RRQ   uint16 = 1
	WRQ   uint16 = 2
	DATA  uint16 = 3
	ACK   uint16 = 4
	ERROR uint16 = 5
	OACK  uint16 = 6
)

type Connection interface {
	ReadFrom(b []byte) (n int, addr net.Addr, err error)
	Write(p []byte) (n int, err error)
}

type RRQresponse struct {
	conn      Connection
	buffer    []byte
	pos       int
	ack       []byte
	blocksize uint16
	blocknum  uint16
}

func (res *RRQresponse) Write(p []byte) (int, error) {
	out := res.buffer[4:]

	bytecount := res.pos + len(p)

	if bytecount < int(res.blocksize) {
		res.pos += copy(out[res.pos:], p)
		return len(p), nil
	}

	copied := copy(out[res.pos:], p)
	remaining := p[copied:]

	res.pos += copied

	if res.pos == int(res.blocksize) {
		_, err := res.writeBuffer()
		if err != nil {
			return 0, err
		}

		if len(remaining) != 0 {
			res.pos = 0
			return res.Write(remaining)
		} else {
			res.pos = 0
		}
	}

	return len(p), nil
}

func (res *RRQresponse) writeBuffer() (int, error) {
	res.blocknum++
	binary.BigEndian.PutUint16(res.buffer, DATA)
	binary.BigEndian.PutUint16(res.buffer[2:], res.blocknum)

	out := res.buffer[:res.pos+4]
	written, err := res.conn.Write(out)

	if err != nil {
		return 0, err
	}
	_, _, err = res.conn.ReadFrom(res.ack)
	if err != nil {
		return 0, err
	}

	// TODO: assert ack

	return written, nil

}

func (res *RRQresponse) End() (int, error) {
	return res.writeBuffer()
}

func NewRRQresponse(addr *net.UDPAddr, blocksize uint16) (*RRQresponse, error) {
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &RRQresponse{
		conn,
		make([]byte, blocksize+4),
		0,
		make([]byte, 5),
		blocksize,
		0,
	}, nil

}

func SendFile(path string, addr *net.UDPAddr) {

	rrq, err := NewRRQresponse(addr, 10)
	if err != nil {
		fmt.Println("Failed to create rrq", err)
		return
	}

	rrq.Write([]byte("Hello "))
	rrq.Write([]byte("Hello "))
	rrq.Write([]byte("Hello "))
	rrq.Write([]byte("Hello "))
	rrq.Write([]byte("Hello "))

	rrq.End()

}

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":1234")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := make([]byte, 20)

	for {
		_, client_addr, err := conn.ReadFrom(data)
		client_udpaddr, err := net.ResolveUDPAddr("udp", client_addr.String())

		if err != nil {
			fmt.Println(err)
			return
		}

		buf := bytes.NewBuffer(data)
		var optcode uint16
		binary.Read(buf, binary.BigEndian, &optcode)
		if optcode == RRQ {
			go SendFile("foo", client_udpaddr)
		} else {
			fmt.Println("got something else", data)
		}

	}

}
