package datatypes

import (
	"bytes"
	"io"
)

type BytesReader interface {
	io.Reader
	io.ByteReader
}

type BytesWriter interface {
	io.Writer
	io.ByteWriter
}

type String struct {
	Size VarInt
	Data string
}

func (s *String) Read(r BytesReader) error {
	var length VarInt
	err := length.Read(r)
	if err != nil {
		return err
	}

	data := make([]byte, length)
	_, err = r.Read(data)
	if err != nil {
		return err
	}

	s.Data = string(data)
	s.Size = length

	return nil
}

func (s *String) Write(w BytesWriter) error {
	err := s.Size.Write(w)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(s.Data))
	if err != nil {
		return err
	}

	return nil
}

func (s *String) String() string {
	return s.Data
}

func (s *String) Length() int {
	return len(s.Data)
}

func (s *String) FromString(str string) {
	s.Size = VarInt(len(str))
	s.Data = str
}

// FromString
// Deprecated
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

// ToString
// Deprecated
func ToString(s String) string {
	return string(s.Data)
}

// WriteString
// Deprecated
func WriteString(s String) []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.Write(WriteVarInt(s.Size))
	buf.Write([]byte(s.Data))
	return buf.Bytes()
}

// ReadString
// Deprecated
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

// ReadStringReader
// Deprecated
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
		Size: VarInt(l),
		Data: string(data),
	}
}
