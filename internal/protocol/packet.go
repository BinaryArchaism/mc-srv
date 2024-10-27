package protocol

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BinaryArchaism/mc-srv/internal/datatypes"
	"github.com/google/uuid"
	"io"
	"os"
)

type PacketWorker struct {
}

func NewPacketWorker() *PacketWorker {
	return &PacketWorker{}
}

type HandshakePacket struct {
	Packet

	ProtocolVersion int
	ServerAddress   datatypes.String
	ServerPort      datatypes.UShort
	NextState       int
}

func (p *HandshakePacket) Read(r io.Reader) error {
	b := Pool.GetN(1024)
	defer Pool.Put(b)

	_, err := r.Read(b)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(b)

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

	p.ServerAddress = datatypes.ReadStringReader(buf)

	err = p.ServerPort.Read(buf)
	if err != nil {
		return err
	}

	p.NextState, err = datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}

	return nil
}

func ReadPacket(r io.Reader) (p []byte, err error) {
	lengthBytes := make([]byte, 3)
	_, err = r.Read(lengthBytes)
	if err != nil {
		return nil, err
	}
	packetLen, l, err := datatypes.ReadVarIntN(lengthBytes)
	if err != nil {
		return nil, err
	}
	res := make([]byte, 0, 3-l+int(packetLen))
	if 3-l != 0 {
		res = append(res, lengthBytes[l:]...)
	}
	_, err = r.Read(res[len(res):packetLen])
	if err != nil {
		return nil, err
	}
	res = res[:packetLen]
	return res, nil
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
	data, err := os.ReadFile("resources/icon.png")
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
	Length int32
	ID     int32

	Name       datatypes.String
	PlayerUUID uuid.UUID
}

func (p *LoginPacket) Read(r io.Reader) error {
	packetBytes := make([]byte, 32)
	_, err := r.Read(packetBytes)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(packetBytes)
	packetLen, err := datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}
	p.Length = int32(packetLen)

	packetID, err := datatypes.BinaryReadVarInt(buf)
	if err != nil {
		return err
	}
	p.ID = int32(packetID)

	name := datatypes.ReadStringReader(buf)
	p.Name = name

	var playerID uuid.UUID
	_, err = buf.Read(playerID[0:])
	if err != nil {
		return err
	}
	p.PlayerUUID = playerID

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

	p.ID = 0x02

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
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

	to, err := resBuf.WriteTo(w)
	if err != nil {
		return err
	}
	fmt.Println(to, p.Length)

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
