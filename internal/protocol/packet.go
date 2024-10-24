package protocol

import (
	"bytes"
	"fmt"
	"github.com/BinaryArchaism/mc-srv/internal/datatypes"
	"io"
)

type HandshakePacket struct {
	Length datatypes.VarInt
	ID     datatypes.VarInt

	ProtocolVersion datatypes.VarInt
	ServerAddress   datatypes.String
	ServerPort      datatypes.Uint16
	NextState       datatypes.VarInt
}

func (p HandshakePacket) WriteTo(w io.Writer) (n int64, err error) {
	buf := bytes.NewBuffer([]byte{})
	buf.Write(datatypes.WriteVarInt(p.ProtocolVersion))
	buf.Write(datatypes.WriteString(p.ServerAddress))
	buf.Write(datatypes.WriteUint16(p.ServerPort))
	buf.Write(datatypes.WriteVarInt(p.NextState))
	res := bytes.NewBuffer([]byte{})
	//p.Length = datatypes.VarInt(buf.Len() + 1)
	p.Length = datatypes.VarInt(0)
	res.Write(datatypes.WriteVarInt(p.Length))
	res.Write(datatypes.WriteVarInt(0))
	res.Write(buf.Bytes())

	for _, b := range res.Bytes() {
		fmt.Printf("%02X", b)
	}
	fmt.Println()

	return res.WriteTo(w)
}
