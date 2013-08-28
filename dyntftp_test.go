package main

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
	return len(p), nil, nil
}

func newRRQResonponse() (*RRQresponse, *MockConnection) {
	conn := &MockConnection{}
	blocksize := 5
	rrq := &RRQresponse{
		conn,
		make([]byte, blocksize+4),
		0,
		make([]byte, 5),
		uint16(blocksize),
		0,
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
	rrq.Write([]byte{1, 2, 3, 4, 5, 6})
	rrq.Write([]byte{7, 8, 9, 0, 0, 9})


	if len(conn.datawritten) != 2 {
		t.Fatalf("Bad value written %v", conn.datawritten)
	}

	if !reflect.DeepEqual(conn.datawritten[1], []byte{0, 3, 0, 2, 6, 7, 8, 9, 0}) {
		t.Fatalf("Bad block written %v", conn.datawritten[1])
	}

	blocknum := binary.BigEndian.Uint16(conn.datawritten[0][2:])
	if blocknum != 1 {
		t.Fatalf("Bad blocknum %v", conn.datawritten)
	}
}
