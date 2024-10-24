package server

import (
	"bufio"
	"context"
	"fmt"
	"github.com/BinaryArchaism/mc-srv/internal/datatypes"
	"github.com/BinaryArchaism/mc-srv/internal/protocol"
	"net"
)

type Server struct {
	srv net.Listener
}

func New() (*Server, error) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return nil, err
	}
	return &Server{
		srv: listener,
	}, nil
}

func (s *Server) Accept(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			err := s.srv.Close()
			if err != nil {
				return err
			}
			return ctx.Err()

		default:
			conn, err := s.srv.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
			}
			go handleConnection(conn)
		}
	}
}

const version = 769

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	fmt.Println("Client connected:", conn.RemoteAddr().String())

	packet := protocol.HandshakePacket{
		ProtocolVersion: version,
		ServerAddress:   datatypes.FromString("localhost"),
		ServerPort:      8080,
		NextState:       2,
	}
	to, err := packet.WriteTo(conn)
	if err != nil {
		fmt.Println("Error writing handshake packet:", err)
		return
	}
	fmt.Println("Handshake packet sent:", to)
	{
		b, err := bufio.NewReader(conn).ReadByte()
		for ; err == nil; b, err = bufio.NewReader(conn).ReadByte() {
			fmt.Printf("ReadByte packet received: %02x\n", b)
		}
	}

	b, err := bufio.NewReader(conn).Peek(1000)
	if err != nil {
		fmt.Println("Error reading handshake packet:", err)
	}
	for _, b := range b {
		fmt.Printf("%02X", b)
	}
	fmt.Println()

	b = []byte{
		0xff, 0x00, 0x23, 0x00, 0xa7, 0x00, 0x31, 0x00, 0x00, 0x00, 0x34, 0x00, 0x37, 0x00, 0x00, 0x00,
		0x31, 0x00, 0x2e, 0x00, 0x34, 0x00, 0x2e, 0x00, 0x32, 0x00, 0x00, 0x00, 0x41, 0x00, 0x20, 0x00,
		0x4f, 0x00, 0x69, 0x00, 0x6e, 0x00, 0x65, 0x00, 0x63, 0x00, 0x72, 0x00, 0x61, 0x00, 0x66, 0x00,
		0x74, 0x00, 0x20, 0x00, 0x53, 0x00, 0x65, 0x00, 0x72, 0x00, 0x76, 0x00, 0x65, 0x00, 0x72, 0x00,
		0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x32, 0x00, 0x30,
	}

	write, err := conn.Write(b)
	if err != nil {
		fmt.Println("Error writing handshake packet:", err)
		return
	}
	fmt.Println("Pong Response packet sent:", write)

	//packet := protocol.HandshakePacket{
	//	ProtocolVersion: version,
	//	ServerAddress:   datatypes.FromString("localhost"),
	//	ServerPort:      8080,
	//	NextState:       2,
	//}
	//to, err := packet.WriteTo(conn)
	//if err != nil {
	//	fmt.Println("Error writing handshake packet:", err)
	//	return
	//}
	//fmt.Println("Handshake packet sent:", to)
	//
	//b, err := bufio.NewReader(conn).Peek(1000)
	//if err != nil {
	//	fmt.Println("Error reading handshake packet:", err)
	//}
	//for _, b := range b {
	//	fmt.Printf("%02X", b)
	//}
	//fmt.Println()
	//
	//write, err := conn.Write([]byte{0xFE, 0x01, 0xFA})
	//if err != nil {
	//	fmt.Println("Error writing handshake packet:", err)
	//	return
	//}
	//fmt.Println("Pong Response packet sent:", write)
	//
	//b, err = bufio.NewReader(conn).Peek(1000)
	//if err != nil {
	//	fmt.Println("Error reading handshake packet:", err)
	//}
	//for _, b := range b {
	//	fmt.Printf("%02X", b)
	//}
	//fmt.Println()
	fmt.Println("==============")
}
