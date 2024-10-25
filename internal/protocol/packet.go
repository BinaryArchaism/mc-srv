package protocol

import (
	"encoding/json"
	"github.com/BinaryArchaism/mc-srv/internal/datatypes"
	"io"
)

// TODO pool of bytes

type HandshakePacket struct {
	Length datatypes.VarInt
	ID     datatypes.VarInt

	ProtocolVersion datatypes.VarInt
	ServerAddress   datatypes.String
	ServerPort      datatypes.UShort
	NextState       datatypes.VarInt
}

func (p *HandshakePacket) Read(r io.Reader) error {
	packetBytes := make([]byte, 32)
	_, err := r.Read(packetBytes)
	if err != nil {
		return err
	}

	packetLen, l, err := datatypes.ReadVarIntN(packetBytes)
	if err != nil {
		return err
	}
	p.Length = packetLen
	packetBytes = packetBytes[l:]

	packetID, l, err := datatypes.ReadVarIntN(packetBytes)
	if err != nil {
		return err
	}
	p.ID = packetID
	packetBytes = packetBytes[l:]

	protocolVersion, l, err := datatypes.ReadVarIntN(packetBytes)
	if err != nil {
		return err
	}
	p.ProtocolVersion = protocolVersion
	packetBytes = packetBytes[l:]

	serverAddress, l := datatypes.ReadStringN(packetBytes)
	p.ServerAddress = serverAddress
	packetBytes = packetBytes[l:]

	serverPort := datatypes.ReadUShort(packetBytes[:2])
	p.ServerPort = serverPort
	packetBytes = packetBytes[2:]

	nextState, l, err := datatypes.ReadVarIntN(packetBytes)
	if err != nil {
		return err
	}
	p.NextState = nextState

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
		Favicon:            image,
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

var image = "data:image/png;base64,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"
