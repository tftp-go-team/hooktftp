
package main

import (
	"net"
	"fmt"
	"encoding/binary"
	"strconv"
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
	blocksize int
	blocknum  uint16
}

func (res *RRQresponse) Write(p []byte) (int, error) {
	out := res.buffer[4:]
	fmt.Println("writing", string(p))

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

		res.pos = 0
		if len(remaining) != 0 {
			return res.Write(remaining)
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
		fmt.Println("failed to write to connection", err)
		return 0, err
	}

	fmt.Println()
	fmt.Println("waiting for ack:", string(out))
	fmt.Println("raw:", (out))
	fmt.Println("size:", len(out))

	_, _, err = res.conn.ReadFrom(res.ack)
	fmt.Println("got ack", res.ack)
	fmt.Println()

	if err != nil {
		fmt.Println("failed to read ack", err)
		return 0, err
	}

	// TODO: assert ack

	return written, nil

}

func (res *RRQresponse) WriteOACK() error {
	if res.blocksize == DEFAULT_BLOCKSIZE {
		return nil
	}

	oackbuffer := make([]byte, 2)
	binary.BigEndian.PutUint16(oackbuffer, OACK)

	oackbuffer = append(oackbuffer, []byte("blksize")...)
	oackbuffer = append(oackbuffer, 0)
	oackbuffer = append(oackbuffer, []byte(strconv.Itoa(res.blocksize))...)
	oackbuffer = append(oackbuffer, 0)

	fmt.Println("oackbuffer", oackbuffer)

	_, err := res.conn.Write(oackbuffer)

	fmt.Println("waiting for oack ack")
	_, _, err = res.conn.ReadFrom(res.ack)
	fmt.Println("got oack", res.ack)
	// TODO: assert ack

	return err
}

func (res *RRQresponse) End() (int, error) {
	return res.writeBuffer()
}

func NewRRQresponse(clientaddr *net.UDPAddr, blocksize int) (*RRQresponse, error) {

	listenaddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", listenaddr, clientaddr)

	if err != nil {
		fmt.Println("failed to create client conn", err)
		return nil, err
	}

	fmt.Println("listenting client on:", conn.LocalAddr(), conn.RemoteAddr())
	return &RRQresponse{
		conn,
		make([]byte, blocksize+4),
		0,
		make([]byte, 4),
		blocksize,
		0,
	}, nil

}

