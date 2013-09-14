package main

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

type RRQrequest struct {
	opcode    uint16
	blocksize int
	mode      int
	path      string
}

type RRQParseError struct {
	msg string
}

func (e *RRQParseError) Error() string {
	return e.msg
}

func sliceUpToNullByte(p []byte) ([]byte, []byte) {
	for i, b := range p {
		if b == 0 {
			return p[0:i], p[i+1 : len(p)]
		}
	}
	return p, nil
}

func ParseRequest(data []byte) (*RRQrequest, error) {
	request := &RRQrequest{blocksize: DEFAULT_BLOCKSIZE}
	request.opcode = binary.BigEndian.Uint16(data)

	if request.opcode != RRQ {
		return request, fmt.Errorf("Unknown optcode %d", request.opcode)
	}

	rest := data[2:len(data)]
	filepath, rest := sliceUpToNullByte(rest)
	request.path = string(filepath)

	mode, rest := sliceUpToNullByte(rest)

	switch string(mode) {
	case "octet":
		request.mode = OCTET
	case "netascii":
		request.mode = NETASCII
	default:
		return request, fmt.Errorf("Unknown mode %v (%v)", mode, string(mode))
	}

	for {
		var option []byte
		option, rest = sliceUpToNullByte(rest)

		if len(option) == 0 {
			break
		}

		switch string(option) {
		case "blksize":
			var blksizebytes []byte
			blksizebytes, rest = sliceUpToNullByte(rest)
			blocksize, err := strconv.Atoi(string(blksizebytes))
			if err != nil {
				fmt.Println("Failed to parse blksize", blksizebytes)
				return request, err
			}
			request.blocksize = blocksize
		default:
			fmt.Println("Unknown option:", option, string(option))
			fmt.Println("data:", data)
			fmt.Println("data string:", string(data))
			// throw away unknown option value
			_, rest = sliceUpToNullByte(rest)
		}

	}

	return request, nil
}
