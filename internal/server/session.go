package server

import (
	"errors"
	"fmt"
	"github.com/BinaryArchaism/mc-srv/internal/datatypes"
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
		s.State = Play
		err = s.PlaySession()
		if err != nil {
			log.Err(err).Msg("failed to play")
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
	var serverBoundPlugin protocol.ServerboundPluginPacket
	err := serverBoundPlugin.Read(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to read serverboundPligin packet: %w", err)
	}

	serverBoundPlugin.Data = nil
	serverBoundPlugin.Channel = datatypes.FromString("")
	err = serverBoundPlugin.Write(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to write serverboundPligin packet: %w", err)
	}

	var clientInfo protocol.ClientInformationPacket
	err = clientInfo.Read(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to read clientInfo packet: %w", err)
	}

	featureFlag := protocol.FeatureFlagPacket{
		TotalFeatures: 0,
		FeatureFlags:  nil,
	}
	err = featureFlag.Write(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to write featureFlag packet: %w", err)
	}

	clientBoundKnownPacksPacket := protocol.KnownPacksPacket{
		KnownPackCount: 1,
		KnownPacks: []protocol.KnownPacks{
			{
				Namespace: datatypes.FromString("minecraft:core"),
				ID:        datatypes.FromString("0"),
				Version:   datatypes.FromString("1.21"),
			},
		},
	}
	err = clientBoundKnownPacksPacket.Write(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to write clientBoundKnownPacksPacket: %w", err)
	}

	var serverBoundKnownPacksPacket protocol.KnownPacksPacket
	err = serverBoundKnownPacksPacket.Read(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to read serverBoundKnownPacksPacket: %w", err)
	}

	// TODO registry data
	// TODO update tags

	finishCfgPacket := protocol.Packet{
		Length: 0x01,
		ID:     0x03,
	}
	err = finishCfgPacket.Write(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to write finishCfgPacker packet: %w", err)
	}

	var clientAckFinishCfgPacket protocol.Packet
	err = clientAckFinishCfgPacket.Read(s.UserConn)
	if err != nil {
		return fmt.Errorf("failed to read clientAckFinishCfgPacket: %w", err)
	}

	return nil
}

func (s *Session) PlaySession() error {
	playLogin := protocol.LoginPlayPacket{
		EntityID:            0,
		IsHardcore:          false,
		DimensionCount:      0,
		DimensionNames:      nil,
		MaxPlayers:          0,
		ViewDistance:        0,
		SimulationDistance:  0,
		ReducedDebugInfo:    false,
		EnableRespawnScreen: false,
		DoLimitedCrafting:   false,
		DimensionType:       0,
		DimensionName:       datatypes.String{},
		HashedSeed:          0,
		GameMode:            0,
		PreviousGameMode:    0,
		IsDebug:             false,
		IsFlat:              false,
		HasDeathLocation:    false,
		DeathDimensionName:  datatypes.String{},
		DeathLocation:       datatypes.Position{},
		PortalCooldown:      0,
		EnforcesSecureChat:  false,
	}
	return nil
}
