package server

import (
	"fmt"
	"github.com/BinaryArchaism/mc-srv/internal/protocol"
	"github.com/rs/zerolog/log"
	"io"
	"net"
)

type State string

const (
	Handshake     State = "handshake"
	Status        State = "status"
	Login         State = "login"
	Configuration State = "configuration"
	Play          State = "play"
)

const (
	statusState = 1
	loginStatus = 2
)

type Session struct {
	State    State
	UserConn io.ReadWriter
}

func NewSession(userConn net.Conn) *Session {
	return &Session{
		UserConn: userConn,
	}
}

func (s *Session) Execute() error {
	s.State = Handshake
	var hsPack protocol.HandshakePacket
	err := hsPack.Read(s.UserConn)
	if err != nil {
		return err
	}
	switch hsPack.NextState {
	case statusState:
		s.State = Status
		err := s.ProcessPingPongSession()
		if err != nil {
			log.Err(err).Msg("failed to process ping pong")
			return err
		}

	case loginStatus:
		s.State = Login

	}
	return nil
}

func (s *Session) ProcessPingPongSession() error {
	var statusRequest protocol.Packet
	err := statusRequest.Read(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to read statusRequest packet: %w", err)
	}

	var statusResponse protocol.StatusResponsePacket
	err = statusResponse.Write(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to write statusResponse packet: %w", err)
	}

	var pingPongPacket protocol.PacketWithData
	err = pingPongPacket.Read(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to read pingPongPacket packet: %w", err)
	}

	err = pingPongPacket.Write(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to write pingPongPacket packet: %w", err)
	}

	return nil
}
