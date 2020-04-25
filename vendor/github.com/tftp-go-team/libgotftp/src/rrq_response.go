package tftp

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type Connection interface {
	SetReadDeadline(time.Time) error
	ReadFrom(b []byte) (n int, addr net.Addr, err error)

	SetWriteDeadline(time.Time) error
	Write(p []byte) (n int, err error)

	Close() (err error)
}

type RRQresponse struct {
	conn         Connection
	buffer       []byte
	pos          int
	ack          []byte
	Request      *Request
	blocknum     uint16
	badinternet  bool
	TransferSize int
}

// readFrom sets a read timeout before calling res.conn.ReadFrom to ensure the
// socket will be free'd if the client is disconnected.
func readFrom(res *RRQresponse, b []byte) (int, net.Addr, error) {
	res.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	return res.conn.ReadFrom(b)
}

// write sets a write timeout before calling res.conn.Write to ensure the
// socket will be free'd if the client is disconnected.
func write(res *RRQresponse, p []byte) (int, error) {
	res.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	return res.conn.Write(p)
}

func (res *RRQresponse) SimulateBadInternet() bool {
	if !res.badinternet {
		return false
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Int()%12345 == 0
}

func (res *RRQresponse) Write(p []byte) (int, error) {
	out := res.buffer[4:]

	bytecount := res.pos + len(p)

	if bytecount < int(res.Request.Blocksize) {
		res.pos += copy(out[res.pos:], p)
		return len(p), nil
	}

	copied := copy(out[res.pos:], p)
	remaining := p[copied:]

	res.pos += copied

	if res.pos == int(res.Request.Blocksize) {

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

	var written int
	if res.SimulateBadInternet() {
		// Just skip sending the packet
		written = len(out)
	} else {
		var err error
		written, err = write(res, out)
		if err != nil {
			return 0, err
		}
	}

	_, _, err := readFrom(res, res.ack)
	if err != nil {
		return 0, err
	}

	opcode := binary.BigEndian.Uint16(res.ack)
	if opcode != ACK {
		return 0, fmt.Errorf("Expected ACK code, got %v", opcode)
	}

	acknum := binary.BigEndian.Uint16(res.ack[2:])
	if acknum == res.blocknum-1 {
		// Got previous ACK <acknum>. Retrying...
		res.blocknum--
		return res.writeBuffer()
	}

	if acknum != res.blocknum {
		return 0, fmt.Errorf(
			"Got weird ACK num %v, expected %v",
			opcode,
			res.blocknum,
		)
	}

	return written, nil
}

func (res *RRQresponse) WriteError(code uint16, message string) error {

	// http://tools.ietf.org/html/rfc1350#page-8
	errorbuffer := make([]byte, 2+2+len(message)+1)

	binary.BigEndian.PutUint16(errorbuffer, ERROR)
	binary.BigEndian.PutUint16(errorbuffer[2:], code)

	copy(errorbuffer[4:], message)
	errorbuffer[len(errorbuffer)-1] = 0

	_, err := write(res, errorbuffer)
	return err
}

func (res *RRQresponse) WriteOACK() error {
	if res.Request.Blocksize == DEFAULT_BLOCKSIZE {
		return nil
	}

	oackbuffer := make([]byte, 2)
	binary.BigEndian.PutUint16(oackbuffer, OACK)

	oackbuffer = append(oackbuffer, []byte("blksize")...)
	oackbuffer = append(oackbuffer, 0)
	oackbuffer = append(oackbuffer, []byte(strconv.Itoa(res.Request.Blocksize))...)
	oackbuffer = append(oackbuffer, 0)

	if res.TransferSize != -1 {
		oackbuffer = append(oackbuffer, []byte("tsize")...)
		oackbuffer = append(oackbuffer, 0)
		oackbuffer = append(oackbuffer, []byte(strconv.Itoa(res.TransferSize))...)
		oackbuffer = append(oackbuffer, 0)
	}

	_, err := write(res, oackbuffer)
	if err != nil {
		return err
	}

	_, _, err = readFrom(res, res.ack)
	if err != nil {
		return err
	}

	// TODO: assert ack

	return nil
}

func (res *RRQresponse) End() (int, error) {
	defer res.conn.Close()

	// Signal end of the transmission. This can be neither empty block or
	// block smaller than res.Request.Blocksize
	return res.writeBuffer()
}

func NewRRQresponse(server *Server, clientaddr *net.UDPAddr, request *Request, badinternet bool) (*RRQresponse, error) {

	listenIp := ""
	if server.listenAddr.IP != nil {
		listenIp = server.listenAddr.IP.String()
	}

	listenaddr, err := net.ResolveUDPAddr("udp", listenIp+":0")
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", listenaddr, clientaddr)
	if err != nil {
		return nil, err
	}

	return &RRQresponse{
		conn,
		make([]byte, request.Blocksize+4),
		0,
		make([]byte, 4),
		request,
		0,
		badinternet,
		-1,
	}, nil

}
