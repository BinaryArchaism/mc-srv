package server

import (
	"errors"
	"fmt"
	"github.com/BinaryArchaism/mc-srv/internal/protocol"
	"github.com/rs/zerolog/log"
	"io"
	"net"
)

var (
	ErrInvalidNextState = errors.New("invalid next state")
	ErrFailedLogin      = errors.New("failed login")
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
		err := s.LoginSession()
		if err != nil {
			log.Err(err).Msg("failed to login")
			return err
		}
		s.State = Configuration
		err = s.ConfigurationSession()
		if err != nil {
			log.Err(err).Msg("failed to configuration")
			return err
		}

	default:
		return ErrInvalidNextState
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

func (s *Session) LoginSession() error {
	var loginPacket protocol.LoginPacket
	err := loginPacket.Read(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to read loginPacket packet: %w", err)
	}

	// TODO encryption
	// encryption skipped

	// TODO set compression
	// compression skipped

	// TODO check player availability to login
	// player wont be disconnected

	loginSuccess := protocol.LoginSuccessPacket{
		UUID:     loginPacket.PlayerUUID,
		UserName: loginPacket.Name,
	}
	err = loginSuccess.Write(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to write loginSuccess packet: %w", err)
	}

	const loginProtocolID = 3
	var loginAck protocol.Packet
	err = loginAck.Read(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to read loginAck packet: %w", err)
	}

	if loginAck.ID != loginProtocolID {
		return ErrFailedLogin
	}

	return nil
}

func (s *Session) ConfigurationSession() error {
	var serverboundPligin protocol.ServerboundPluginPacket
	err := serverboundPligin.Read(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to read serverboundPligin packet: %w", err)
	}

	err = serverboundPligin.Write(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to write serverboundPligin packet: %w", err)
	}

	return nil
}
