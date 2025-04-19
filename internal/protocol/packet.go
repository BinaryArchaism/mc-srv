package protocol

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/BinaryArchaism/mc-srv/internal/countingbuffer"
	"github.com/BinaryArchaism/mc-srv/internal/datatypes"
	"github.com/google/uuid"
)

type HandshakePacket struct {
	Packet

	ProtocolVersion int
	ServerAddress   string
	ServerPort      int
	NextState       int
}

func (p *HandshakePacket) Read(r io.Reader) error {
	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	_, err := r.Read(poolBytes[:cap(poolBytes)])
	if err != nil {
		return err
	}

	buf := countingbuffer.New(poolBytes)

	p.Length, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	p.ID, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	p.ProtocolVersion, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	serverAddress := datatypes.ReadStringReader(buf)
	p.ServerAddress = serverAddress.Data

	var serverPort datatypes.UShort
	err = serverPort.Read(buf)
	if err != nil {
		return err
	}
	p.ServerPort = int(serverPort)

	p.NextState, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	return nil
}

type StatusResponsePacket struct {
	Length datatypes.VarInt
	ID     datatypes.VarInt

	JSONResponse datatypes.String
}

type JSONResponse struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description struct {
		Text string `json:"text"`
	} `json:"description"`
	Favicon            string `json:"favicon"`
	EnforcesSecureChat bool   `json:"enforcesSecureChat"`
}

func (p *StatusResponsePacket) Write(w io.Writer) error {
	// todo must do it once on startup
	data, err := os.ReadFile("../../resources/icon.png")
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(data)

	js := JSONResponse{
		Version: struct {
			Name     string `json:"name"`
			Protocol int    `json:"protocol"`
		}{
			Name:     "1.21",
			Protocol: 767,
		},
		Players: struct {
			Max    int `json:"max"`
			Online int `json:"online"`
			Sample []struct {
				Name string `json:"name"`
				Id   string `json:"id"`
			} `json:"sample"`
		}{
			Max:    100,
			Online: 10,
			Sample: nil,
		},
		Description: struct {
			Text string `json:"text"`
		}{
			"Davai rabotai",
		},
		Favicon:            "data:image/png;base64," + encoded,
		EnforcesSecureChat: false,
	}
	b, err := json.Marshal(js)
	if err != nil {
		return err
	}
	str := datatypes.FromString(string(b))

	strBytes := datatypes.WriteString(str)

	s := StatusResponsePacket{
		Length: datatypes.VarInt(len(strBytes) + 1),
		ID:     0x00,
	}

	res := make([]byte, 0, 3+len(b)+1)
	res = append(res, datatypes.WriteVarInt(s.Length)...)
	res = append(res, datatypes.WriteVarInt(s.ID)...)
	res = append(res, strBytes...)

	_, err = w.Write(res)
	if err != nil {
		return err
	}
	return nil
}

type LoginPacket struct {
	Packet

	Name       datatypes.String
	PlayerUUID uuid.UUID
}

func (p *LoginPacket) Read(r io.Reader) error {
	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	_, err := r.Read(poolBytes[:cap(poolBytes)])
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(poolBytes)

	p.Length, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	p.ID, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	p.Name = datatypes.ReadStringReader(buf)

	_, err = buf.Read(p.PlayerUUID[:])
	if err != nil {
		return err
	}

	return nil
}

type DisconnectPacket struct {
	Length int32
	ID     int32

	Reason datatypes.String
}

func (p *DisconnectPacket) Write(w io.Writer) error {
	type jsonText struct {
		Text string `json:"text"`
	}
	js := jsonText{
		Text: "restricted",
	}
	b, err := json.Marshal(js)
	if err != nil {
		return err
	}
	str := datatypes.FromString(string(b))
	strBytes := datatypes.WriteString(str)

	s := DisconnectPacket{
		Length: int32(datatypes.VarInt(len(strBytes) + 1)),
		ID:     0x00,
	}

	res := make([]byte, 0, 3+len(b)+1)
	res = append(res, datatypes.WriteVarInt(datatypes.VarInt(s.Length))...)
	res = append(res, datatypes.WriteVarInt(datatypes.VarInt(s.ID))...)
	res = append(res, strBytes...)

	_, err = w.Write(res)
	if err != nil {
		return err
	}
	return nil
}

type Packet struct {
	Length int
	ID     int
}

func (p *Packet) Read(r io.Reader) error {
	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	_, err := r.Read(poolBytes[:cap(poolBytes)])
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(poolBytes)

	p.Length, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}
	p.ID, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	return nil
}

func (p *Packet) Write(w io.Writer) error {
	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	buf := bytes.NewBuffer(poolBytes)
	buf.Reset()

	_, err := buf.Write(datatypes.BinaryWriteVarInt(p.Length))
	if err != nil {
		return err
	}
	_, err = buf.Write(datatypes.BinaryWriteVarInt(p.ID))
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}
	return nil
}

type PacketWithData struct {
	Packet
	Data []byte
}

func (p *PacketWithData) Read(r io.Reader) error {
	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	_, err := r.Read(poolBytes[:cap(poolBytes)])
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(poolBytes)

	p.Length, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}
	p.ID, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	p.Data = poolBytes[:p.Length-1]
	return nil
}

func (p *PacketWithData) Write(w io.Writer) error {
	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	buf := bytes.NewBuffer(poolBytes)
	buf.Reset()

	_, err := buf.Write(datatypes.BinaryWriteVarInt(p.Length))
	if err != nil {
		return err
	}
	_, err = buf.Write(datatypes.BinaryWriteVarInt(p.ID))
	if err != nil {
		return err
	}

	_, err = buf.Write(p.Data)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

type LoginSuccessPacket struct {
	Packet

	UUID                uuid.UUID
	UserName            datatypes.String
	NumOfProps          int
	Property            []Property
	StrictErrorHandling datatypes.Boolean
}

type Property struct {
	Name      datatypes.String
	Value     datatypes.String
	IsSigned  datatypes.Boolean
	Signature datatypes.String
}

func (p *LoginSuccessPacket) Write(w io.Writer) error {
	if p.NumOfProps != len(p.Property) {
		return errors.New("invalid number of props")
	}

	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	buf := bytes.NewBuffer(poolBytes)
	buf.Reset()

	p.ID = 0x02

	buf.Write(datatypes.WriteVarInt(datatypes.VarInt(p.ID)))
	buf.Write(p.UUID[:])
	buf.Write(datatypes.WriteString(p.UserName))
	buf.Write(datatypes.WriteVarInt(datatypes.VarInt(p.NumOfProps)))

	for i := range p.NumOfProps {
		buf.Write(datatypes.WriteString(p.Property[i].Name))
		buf.Write(datatypes.WriteString(p.Property[i].Value))
		buf.WriteByte(datatypes.WriteBoolean(p.Property[i].IsSigned))
		if p.Property[i].IsSigned {
			buf.Write(datatypes.WriteString(p.Property[i].Signature))
		}
	}
	buf.WriteByte(datatypes.WriteBoolean(p.StrictErrorHandling))

	p.Length = buf.Len()

	resBuf := bytes.Buffer{}
	resBuf.Write(datatypes.WriteVarInt(datatypes.VarInt(p.Length)))
	resBuf.Write(buf.Bytes())

	_, err := resBuf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

type ClientboundKnownPacksPacket struct {
	Packet

	KnownPacketCount int
	KnownPacket      []KnownPack
}

type KnownPack struct {
	Namespace datatypes.String
	ID        datatypes.String
	Version   datatypes.String
}

func (p *ClientboundKnownPacksPacket) Write(w io.Writer) error {
	if p.KnownPacketCount != len(p.KnownPacket) {
		return errors.New("invalid number of props")
	}

	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	buf := bytes.NewBuffer(poolBytes)
	buf.Reset()

	p.Packet.ID = 0x0E
	buf.Write(datatypes.BinaryWriteVarInt(p.ID))
	buf.Write(datatypes.BinaryWriteVarInt(p.KnownPacketCount))
	for _, kp := range p.KnownPacket {
		buf.Write(datatypes.WriteString(kp.Namespace))
		buf.Write(datatypes.WriteString(kp.ID))
		buf.Write(datatypes.WriteString(kp.Version))
	}

	p.Length = buf.Len() + 1

	l, err := w.Write(datatypes.BinaryWriteVarInt(p.Length))
	if err != nil {
		return err
	}

	n, err := buf.WriteTo(w)
	if err != nil {
		return err
	}

	if n+int64(l) != int64(p.Length) {
		return io.ErrShortWrite
	}

	return nil
}

func (p *ClientboundKnownPacksPacket) Read(r io.Reader) error {
	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	_, err := r.Read(poolBytes)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(poolBytes[:])

	packetLen, err := datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}
	packetID, err := datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	knownPackCount, err := datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	knownPacks := make([]KnownPack, 0, knownPackCount)
	var tmpKP KnownPack
	for i := 0; i < int(knownPackCount); i++ {
		tmpKP.Namespace = datatypes.ReadStringReader(buf)
		tmpKP.ID = datatypes.ReadStringReader(buf)
		tmpKP.Version = datatypes.ReadStringReader(buf)
		knownPacks = append(knownPacks, tmpKP)
	}

	p.Packet.Length = int(packetLen)
	p.Packet.ID = int(packetID)

	p.KnownPacketCount = int(knownPackCount)
	p.KnownPacket = knownPacks

	return nil
}

type LoginPlugin struct {
	Packet

	MessageID  int
	Successful datatypes.Boolean
	Data       []byte
}

func (p *LoginPlugin) Read(r io.Reader) error {
	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)
	_, err := r.Read(poolBytes)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(poolBytes)

	p.Packet.Length, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}
	p.Packet.ID, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}
	p.MessageID, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}
	p.Data = buf.Bytes()[:p.Length-2]

	return nil
}

type ServerboundPluginPacket struct {
	Packet

	Channel datatypes.String
	Data    []byte
}

func (p *ServerboundPluginPacket) Read(r io.Reader) error {
	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	n, err := r.Read(poolBytes[:cap(poolBytes)])
	if err != nil {
		return err
	}

	buf := countingbuffer.New(poolBytes)

	p.Packet.Length, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	p.Packet.ID, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	p.Channel = datatypes.ReadStringReader(buf)
	p.Data = buf.Bytes()[:n-buf.ReadCount()]
	return nil
}

func (p *ServerboundPluginPacket) Write(w io.Writer) error {
	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	buf := countingbuffer.New(poolBytes)
	buf.Reset()

	p.ID = 0x01

	_, err := buf.Write(datatypes.BinaryWriteVarInt(p.ID))
	if err != nil {
		return err
	}
	_, err = buf.Write(datatypes.WriteString(p.Channel))
	if err != nil {
		return err
	}

	_, err = w.Write(datatypes.BinaryWriteVarInt(buf.WriteCount()))
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}
	return nil
}

type ClientInformationPacket struct {
	Packet

	Locale              datatypes.String
	ViewDistance        byte
	ChatMode            int
	ChatColors          datatypes.Boolean
	DisplayedSkinParts  byte
	MainHand            int
	EnableTextFiltering datatypes.Boolean
	AllowServerListings datatypes.Boolean
}

func (p *ClientInformationPacket) Read(r io.Reader) error {
	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	_, err := r.Read(poolBytes[:cap(poolBytes)])
	if err != nil {
		return err
	}

	buf := countingbuffer.New(poolBytes)

	p.Packet.Length, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	p.Packet.ID, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	p.Locale = datatypes.ReadStringReader(buf)
	p.ViewDistance, err = buf.ReadByte()
	if err != nil {
		return err
	}

	p.ChatMode, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	chatColors, err := buf.ReadByte()
	if err != nil {
		return err
	}
	p.ChatColors = datatypes.ReadBoolean(chatColors)

	p.DisplayedSkinParts, err = buf.ReadByte()
	if err != nil {
		return err
	}

	p.MainHand, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	enableTextFiltering, err := buf.ReadByte()
	if err != nil {
		return err
	}
	p.EnableTextFiltering = datatypes.ReadBoolean(enableTextFiltering)

	allowServerListings, err := buf.ReadByte()
	if err != nil {
		return err
	}
	p.AllowServerListings = datatypes.ReadBoolean(allowServerListings)

	return nil
}

type FeatureFlagPacket struct {
	Packet

	TotalFeatures int
	FeatureFlags  []datatypes.String
}

func (p *FeatureFlagPacket) Write(w io.Writer) error {
	if p.TotalFeatures != len(p.FeatureFlags) {
		return errors.New("invalid number of FeatureFlags")
	}

	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	buf := countingbuffer.New(poolBytes)
	buf.Reset()

	p.ID = 0x0C

	_, err := buf.Write(datatypes.BinaryWriteVarInt(p.ID))
	if err != nil {
		return err
	}
	_, err = buf.Write(datatypes.BinaryWriteVarInt(p.TotalFeatures))
	if err != nil {
		return err
	}

	for _, f := range p.FeatureFlags {
		_, err = buf.Write(datatypes.WriteString(f))
		if err != nil {
			return err
		}
	}

	_, err = w.Write(datatypes.BinaryWriteVarInt(buf.WriteCount()))
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}
	return nil
}

type KnownPacksPacket struct {
	Packet

	KnownPackCount int
	KnownPacks     []KnownPacks
}

type KnownPacks struct {
	Namespace datatypes.String
	ID        datatypes.String
	Version   datatypes.String
}

func (p *KnownPacksPacket) Write(w io.Writer) error {
	if p.KnownPackCount != len(p.KnownPacks) {
		return errors.New("invalid number of KnownPacks")
	}

	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	buf := countingbuffer.New(poolBytes)
	buf.Reset()

	p.ID = 0x0E

	_, err := buf.Write(datatypes.BinaryWriteVarInt(p.ID))
	if err != nil {
		return err
	}
	_, err = buf.Write(datatypes.BinaryWriteVarInt(p.KnownPackCount))
	if err != nil {
		return err
	}

	for _, kp := range p.KnownPacks {
		_, err = buf.Write(datatypes.WriteString(kp.Namespace))
		if err != nil {
			return err
		}
		_, err = buf.Write(datatypes.WriteString(kp.ID))
		if err != nil {
			return err
		}
		_, err = buf.Write(datatypes.WriteString(kp.Version))
		if err != nil {
			return err
		}
	}

	_, err = w.Write(datatypes.BinaryWriteVarInt(buf.WriteCount()))
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}
	return nil
}

func (p *KnownPacksPacket) Read(r io.Reader) error {
	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	_, err := r.Read(poolBytes[:cap(poolBytes)])
	if err != nil {
		return err
	}

	buf := countingbuffer.New(poolBytes)

	p.Packet.Length, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	p.Packet.ID, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	p.KnownPackCount, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	for range p.KnownPackCount {
		p.KnownPacks = append(p.KnownPacks, KnownPacks{
			Namespace: datatypes.ReadStringReader(buf),
			ID:        datatypes.ReadStringReader(buf),
			Version:   datatypes.ReadStringReader(buf),
		})
	}

	return nil
}

type LoginPlayPacket struct {
	Packet

	EntityID            int32
	IsHardcore          datatypes.Boolean
	DimensionCount      datatypes.VarInt
	DimensionNames      []datatypes.String
	MaxPlayers          datatypes.VarInt
	ViewDistance        datatypes.VarInt
	SimulationDistance  datatypes.VarInt
	ReducedDebugInfo    datatypes.Boolean
	EnableRespawnScreen datatypes.Boolean
	DoLimitedCrafting   datatypes.Boolean
	DimensionType       datatypes.VarInt
	DimensionName       datatypes.String
	HashedSeed          int64
	GameMode            byte
	PreviousGameMode    byte
	IsDebug             datatypes.Boolean
	IsFlat              datatypes.Boolean
	HasDeathLocation    datatypes.Boolean
	DeathDimensionName  datatypes.String
	DeathLocation       datatypes.Position
	PortalCooldown      datatypes.VarInt
	EnforcesSecureChat  datatypes.Boolean
}

func (p *LoginPlayPacket) Write(w io.Writer) error {
	if int(p.DimensionCount) != len(p.DimensionNames) {
		return errors.New("invalid number of KnownPacks")
	}

	poolBytes := Pool.GetN(SmallObjectSize)
	defer Pool.Put(poolBytes)

	buf := countingbuffer.New(poolBytes)
	buf.Reset()

	p.ID = 0x2B

	binary.BigEndian.PutUint32(buf.Next(4), uint32(p.EntityID))

	err := buf.WriteByte(datatypes.WriteBoolean(p.IsHardcore))
	if err != nil {
		return err
	}

	_, err = buf.Write(datatypes.BinaryWriteVarInt(int(p.DimensionCount)))
	if err != nil {
		return err
	}

	for _, d := range p.DimensionNames {
		_, err = buf.Write(datatypes.WriteString(d))
		if err != nil {
			return err
		}
	}

	_, err = buf.Write(datatypes.BinaryWriteVarInt(int(p.MaxPlayers)))
	if err != nil {
		return err
	}

	_, err = buf.Write(datatypes.BinaryWriteVarInt(int(p.ViewDistance)))
	if err != nil {
		return err
	}

	_, err = buf.Write(datatypes.BinaryWriteVarInt(int(p.SimulationDistance)))
	if err != nil {
		return err
	}

	err = buf.WriteByte(datatypes.WriteBoolean(p.ReducedDebugInfo))
	if err != nil {
		return err
	}

	err = buf.WriteByte(datatypes.WriteBoolean(p.EnableRespawnScreen))
	if err != nil {
		return err
	}

	err = buf.WriteByte(datatypes.WriteBoolean(p.DoLimitedCrafting))
	if err != nil {
		return err
	}

	_, err = buf.Write(datatypes.BinaryWriteVarInt(int(p.DimensionType)))
	if err != nil {
		return err
	}

	_, err = buf.Write(datatypes.WriteString(p.DimensionName))
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint64(buf.Next(8), uint64(p.HashedSeed))

	err = buf.WriteByte(p.GameMode)
	if err != nil {
		return err
	}

	err = buf.WriteByte(p.PreviousGameMode)
	if err != nil {
		return err
	}

	err = buf.WriteByte(datatypes.WriteBoolean(p.IsDebug))
	if err != nil {
		return err
	}

	err = buf.WriteByte(datatypes.WriteBoolean(p.IsFlat))
	if err != nil {
		return err
	}

	err = buf.WriteByte(datatypes.WriteBoolean(p.HasDeathLocation))
	if err != nil {
		return err
	}

	if p.HasDeathLocation {
		_, err = buf.Write(datatypes.WriteString(p.DeathDimensionName))
		if err != nil {
			return err
		}

		binary.LittleEndian.PutUint64(buf.Next(8), uint64(datatypes.WritePosition(p.DeathLocation)))
	}

	_, err = buf.Write(datatypes.BinaryWriteVarInt(int(p.PortalCooldown)))
	if err != nil {
		return err
	}

	err = buf.WriteByte(datatypes.WriteBoolean(p.EnforcesSecureChat))
	if err != nil {
		return err
	}

	_, err = w.Write(datatypes.BinaryWriteVarInt(buf.WriteCount()))
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}
