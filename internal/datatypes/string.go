package datatypes

import (
	"bytes"
	"io"
)

type String struct {
	Size VarInt
	Data string
}

func FromString(s string) String {
	length := VarInt(len(s))
	data := make([]byte, 0, length)
	for _, r := range s {
		data = append(data, byte(r))
	}
	return String{
		Size: length,
		Data: string(data),
	}
}

func ToString(s String) string {
	return string(s.Data)
}

func WriteString(s String) []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.Write(WriteVarInt(s.Size))
	buf.Write([]byte(s.Data))
	return buf.Bytes()
}

func ReadString(b []byte) String {
	length, _ := ReadVarInt(b)
	data := make([]byte, 0, length)
	for i := len(b) - int(length); i < len(b); i++ {
		data = append(data, b[i])
	}
	return String{
		Size: length,
		Data: string(data),
	}
}

func ReadStringN(b []byte) (String, int) {
	length, l, _ := ReadVarIntN(b)
	data := make([]byte, 0, length)
	for i := l; i < int(length)+l; i++ {
		data = append(data, b[i])
	}
	return String{
		Size: length,
		Data: string(data),
	}, len(data) + l
}

func ReadStringReader(in io.ByteReader) String {
	l, err := BinaryReadVarInt(in)
	if err != nil {
		return String{}
	}
	data := make([]byte, 0, l)
	for i := 0; i < int(l); i++ {
		b, err := in.ReadByte()
		if err != nil {
			return String{}
		}
		data = append(data, b)
	}
	return String{
		Size: l,
		Data: string(data),
	}
}

func ReadStrings(byteStrings []byte) []String {
	var res []String
	for i := 0; i < len(byteStrings); i++ {
		length, _ := ReadVarInt(byteStrings)
		varIntBytesLen := len(WriteVarInt(length))
		str := byteStrings[i+varIntBytesLen+1 : i+varIntBytesLen+1+int(length)]
		res = append(res, String{
			Size: length,
			Data: string(str),
		})
	}
	return res
}
