package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
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
	defer conn.Close()
	fmt.Println("Client connected:", conn.RemoteAddr().String())

	for {

		// Read message from the client
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr().String())
			return
		}

		// Print the message
		fmt.Print("Message received:", string(message))

		// Respond back to the client
		newMessage := strings.ToUpper(message)
		conn.Write([]byte(newMessage))
	}
}
