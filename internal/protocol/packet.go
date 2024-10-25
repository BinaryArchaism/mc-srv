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
