package tftp

import (
	"fmt"
	"net"
)

type Server struct {
	conn *net.UDPConn
	buffer []byte
}

func (server *Server) WaitForRequest() (*Request, *RRQresponse, error) {

		written, addr, err := server.conn.ReadFrom(server.buffer)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to read data from client: %v", err)
		}

		request, err := ParseRequest(server.buffer[:written])
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to parse request: %v", err)
		}

		if request.Opcode != RRQ {
			return nil, nil, fmt.Errorf("Unkown opcode %v", request.Opcode)
		}

		raddr, err := net.ResolveUDPAddr("udp", addr.String())
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to resolve client address: %v", err)
		}

		response, err := NewRRQresponse(raddr, request.Blocksize, false)
		if err != nil {
			return nil, nil, err
		}

		return request, response, nil
}


func NewTFTPServer(addr *net.UDPAddr) (*Server, error){

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("Failed listen UDP %v", err)
	}

	return &Server{
		conn,
		make([]byte, 50),
	}, nil

}
