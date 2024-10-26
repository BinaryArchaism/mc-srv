package server

import (
	"context"
	"fmt"
	"github.com/BinaryArchaism/mc-srv/internal/protocol"
	"io"
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
			go s.HandleConnection(conn)
		}
	}
}

const (
	statusState = 1
	loginStatus = 2
)

func (s *Server) HandleConnection(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection:", err)
		}
	}(conn)
	fmt.Println("Client connected:", conn.RemoteAddr().String())

	fmt.Println("=== HANDSHAKE SESSION ===")
	var hsPacket protocol.HandshakePacket
	err := hsPacket.Read(conn)
	if err != nil {
		fmt.Println("Error reading handshake packet:", err)
		return
	}
	fmt.Printf(" <- Client handshake packet received: %+v\n", hsPacket)

	switch hsPacket.NextState {
	case statusState:
		err = s.PingSession(conn)
		if err != nil {
			fmt.Println("Error pinging session:", err)
			return
		}
	case loginStatus:
		err = s.LoginSession(conn)
		if err != nil {
			fmt.Println("Error logging in:", err)
			return
		}
	default:
		fmt.Println("Unknown handshake packet:", hsPacket.NextState)
		return
	}
}

func (s *Server) LoginSession(conn net.Conn) error {
	var loginPacket protocol.LoginPacket
	err := loginPacket.Read(conn)
	if err != nil {
		return err
	}
	//var disconnectPacket protocol.DisconnectPacket
	//err = disconnectPacket.Write(conn)
	//if err != nil {
	//	return err
	//}
	loginSuccessPacket := protocol.LoginSuccessPacket{
		UUID:       loginPacket.PlayerUUID,
		UserName:   loginPacket.Name,
		NumOfProps: 0,
		Property:   nil,
	}
	err = loginSuccessPacket.Write(conn)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) PingSession(conn net.Conn) error {
	fmt.Println("=== STATUS SESSION ===")
	b, err := readAll(conn)
	if err != nil {
		fmt.Println("Error reading packet:", err)
		return err
	}
	fmt.Printf(" <- Client status request: %X\n", b)

	err = (&protocol.StatusResponsePacket{}).Write(conn)
	if err != nil {
		fmt.Println("Error writing packet:", err)
		return err
	}
	fmt.Println(" -> Server status response")

	fmt.Println("=== PING SESSION ===")

	b, err = readAll(conn)
	if err != nil {
		fmt.Println("Error reading packet:", err)
		return err
	}
	if len(b) == 0 {
		fmt.Println("Error reading packet: buffer is empty")
		return err
	}
	fmt.Printf(" <- Client ping request: %X\n", b)
	_, err = conn.Write(b)
	if err != nil {
		fmt.Println("Error writing packet:", err)
		return err
	}
	fmt.Println(" -> Server ping response")
	fmt.Println("=== CLOSE CONNECTION ===")
	return nil
}

func readAll(r io.Reader) ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}

		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		return b, nil
	}
}
