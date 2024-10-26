package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/BinaryArchaism/mc-srv/internal/datatypes"
	"github.com/BinaryArchaism/mc-srv/internal/protocol"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"io"
	"net"
)

type Server struct {
	srv net.Listener
}

func New() (*Server, error) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Err(err).Msg("Error starting TCP server")
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
			log.Err(err).Msg("Error closing connection")
		}
		log.Trace().Str("client", conn.RemoteAddr().String()).Msg("Connection closed")
	}(conn)
	log.Trace().Msgf("Client connected: %s", conn.RemoteAddr().String())

	var hsPacket protocol.HandshakePacket
	err := hsPacket.Read(conn)
	if err != nil {
		log.Err(err).Msg("Error reading handshake packet")
		return
	}
	log.Trace().Any("handshake packet", hsPacket).Msg("serverbound")

	switch hsPacket.NextState {
	case statusState:
		err = s.PingSession(conn)
		if err != nil {
			log.Err(err).Msg("Error while ping session")
			return
		}
	case loginStatus:
		err = s.LoginSession(conn)
		if err != nil {
			log.Err(err).Msg("Error while logging")
			return
		}
	default:
		log.Error().Int("handshake packet next state", int(hsPacket.NextState)).Msg("Unknown handshake packet")
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
		UUID:     uuid.MustParse("4566e69fc90748ee8d71d7ba5aa00d20"),
		UserName: datatypes.FromString("Thinkofdeath"),
	}
	err = loginSuccessPacket.Write(conn)
	if err != nil {
		return err
	}

	b, err := readAll(conn)
	if err != nil {
		return err
	}
	log.Trace().Bytes("login ack packet", b).Msg("serverbound")

	return nil
}

func (s *Server) PingSession(conn net.Conn) error {
	b, err := readAll(conn)
	if err != nil {
		return fmt.Errorf("error reading status request: %w", err)
	}
	log.Trace().Bytes("status request packet", b).Msg("serverbound")

	var statusResponse protocol.StatusResponsePacket
	err = statusResponse.Write(conn)
	if err != nil {
		return fmt.Errorf("error writing status response: %w", err)
	}
	log.Trace().Any("status response packet", statusResponse).Msg("clientbound")

	b, err = readAll(conn)
	if err != nil {
		return fmt.Errorf("error reading ping request: %w", err)
	}
	if len(b) == 0 {
		return errors.New("ping request is empty")
	}
	log.Trace().Bytes("ping request packet", b).Msg("serverbound")

	_, err = conn.Write(b)
	if err != nil {
		return fmt.Errorf("error writing pong response: %w", err)
	}
	log.Trace().Bytes("pong response packet", b).Msg("clientbound")

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
