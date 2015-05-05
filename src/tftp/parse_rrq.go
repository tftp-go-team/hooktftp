package tftp

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"net"
	"logger"
)

type Request struct {
	Opcode    uint16
	Blocksize int
	Mode      int
	Path      string
	Addr      *net.Addr
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

func ParseRequest(data []byte) (*Request, error) {
	request := &Request{Blocksize: DEFAULT_BLOCKSIZE}
	request.Opcode = binary.BigEndian.Uint16(data)

	if request.Opcode != RRQ {
		return request, fmt.Errorf("Unknown optcode %d", request.Opcode)
	}

	rest := data[2:len(data)]
	filepath, rest := sliceUpToNullByte(rest)
	request.Path = string(filepath)

	mode, rest := sliceUpToNullByte(rest)

	switch string(mode) {
	case "octet":
		request.Mode = OCTET
	case "netascii":
		request.Mode = NETASCII
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
				logger.Err("Failed to parse blksize %d", blksizebytes)
				return request, err
			}
			request.Blocksize = blocksize
		default:
			logger.Err("Unknown option: %s; data: %s", string(option), string(data))
			// throw away unknown option value
			_, rest = sliceUpToNullByte(rest)
		}

	}

	return request, nil
}
