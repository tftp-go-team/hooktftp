package tftp

import (
	"encoding/binary"
	"net"
	"reflect"
	"testing"
)

type MockConnection struct {
	datawritten [][]byte
}

func (conn *MockConnection) Write(incoming []byte) (int, error) {
	data := make([]byte, len(incoming))
	copy(data, incoming)
	conn.datawritten = append(conn.datawritten, data)
	return len(incoming), nil
}

func (conn *MockConnection) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	num := binary.BigEndian.Uint16(conn.datawritten[len(conn.datawritten)-1][2:])
	binary.BigEndian.PutUint16(p, ACK)
	binary.BigEndian.PutUint16(p[2:], num)
	return len(p), nil, nil
}

func newRRQResonponse() (*RRQresponse, *MockConnection) {
	conn := &MockConnection{}
	request := &Request{
		RRQ,
		5,
		OCTET,
		-1,
		"/foo",
		nil,
	}
	rrq := &RRQresponse{
		conn,
		make([]byte, request.Blocksize+4),
		0,
		make([]byte, 5),
		request,
		0,
		false,
		-1,
	}
	return rrq, conn
}

func TestSmallWrite(t *testing.T) {
	rrq, conn := newRRQResonponse()
	rrq.Write([]byte{1, 2, 3})
	if len(conn.datawritten) != 0 {
		t.Fatalf("Data written too early")
	}
}

func Test2SmallWrites(t *testing.T) {
	rrq, conn := newRRQResonponse()
	rrq.Write([]byte{1, 2})
	rrq.Write([]byte{3, 4})

	if len(conn.datawritten) != 0 {
		t.Fatalf("Data written too early")
	}
}

func TestFullPacket(t *testing.T) {
	rrq, conn := newRRQResonponse()
	rrq.Write([]byte{1, 2, 3, 4, 5})

	if len(conn.datawritten) != 1 {
		t.Fatalf("Bad value written %v", conn.datawritten)
	}

	blocknum := binary.BigEndian.Uint16(conn.datawritten[0][2:])
	if blocknum != 1 {
		t.Fatalf("Bad blocknum %v", conn.datawritten)
	}

	if !reflect.DeepEqual(conn.datawritten[0], []byte{0, 3, 0, 1, 1, 2, 3, 4, 5}) {
		t.Fatalf("Bad first value written %v", conn.datawritten[0])
	}

}

func Test2FullPackets(t *testing.T) {
	rrq, conn := newRRQResonponse()
	rrq.Write([]byte{1, 2, 3, 4, 5})
	rrq.Write([]byte{1, 2, 3, 4, 5})

	if len(conn.datawritten) != 2 {
		t.Fatalf("Bad value written %v", conn.datawritten)
	}

	blocknum := binary.BigEndian.Uint16(conn.datawritten[0][2:])
	if blocknum != 1 {
		t.Fatalf("Bad blocknum %v", conn.datawritten)
	}

	if !reflect.DeepEqual(conn.datawritten[0], []byte{0, 3, 0, 1, 1, 2, 3, 4, 5}) {
		t.Fatalf("Bad block value written %v", conn.datawritten[0])
	}

	if !reflect.DeepEqual(conn.datawritten[1], []byte{0, 3, 0, 2, 1, 2, 3, 4, 5}) {
		t.Fatalf("Bad block value written %v", conn.datawritten[0])
	}

}

func TestLargePacket(t *testing.T) {
	rrq, conn := newRRQResonponse()
	rrq.Write([]byte{1, 2, 3, 4, 5, 6})

	if len(conn.datawritten) != 1 {
		t.Fatalf("Bad value written %v", conn.datawritten)
	}

	if !reflect.DeepEqual(conn.datawritten[0], []byte{0, 3, 0, 1, 1, 2, 3, 4, 5}) {
		t.Fatalf("Bad first block written %v", conn.datawritten[0])
	}

}

func Test2LargePackets(t *testing.T) {
	rrq, conn := newRRQResonponse()

	_, err := rrq.Write([]byte{1, 2, 3, 4, 5, 6})
	if err != nil {
		t.Fatalf("Got error after write: %v", err)
	}

	_, err = rrq.Write([]byte{7, 8, 9, 0, 0, 9})
	if err != nil {
		t.Fatalf("Got error after write: %v", err)
	}

	if len(conn.datawritten) != 2 {
		t.Fatalf("Bad value written %v", conn.datawritten)
	}

	if !reflect.DeepEqual(conn.datawritten[1], []byte{0, 3, 0, 2, 6, 7, 8, 9, 0}) {
		t.Fatalf("Bad block written g %v", conn.datawritten[1])
	}

	blocknum := binary.BigEndian.Uint16(conn.datawritten[0][2:])
	if blocknum != 1 {
		t.Fatalf("Bad blocknum %v", conn.datawritten)
	}
}

func TestOACK(t *testing.T) {
	rrq, conn := newRRQResonponse()
	rrq.WriteOACK()

	if !reflect.DeepEqual(conn.datawritten[0], []byte{0, 6, 98, 108, 107, 115, 105, 122, 101, 0, 53, 0}) {
		t.Fatalf("Bad oack written %v", conn.datawritten[0])
	}

}

func TestErrorPacket(t *testing.T) {
	rrq, conn := newRRQResonponse()
	// http://tools.ietf.org/html/rfc1350#page-8
	rrq.WriteError(NOT_FOUND, "aa")

	if !reflect.DeepEqual(conn.datawritten[0], []byte{0, 5, 0, 1, 97, 97, 0}) {
		t.Fatalf("Bad error written %v", conn.datawritten[0])
	}

}
