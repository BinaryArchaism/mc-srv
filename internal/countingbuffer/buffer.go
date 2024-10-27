package countingbuffer

import (
	"bytes"
	"io"
)

type CountingBuffer struct {
	*bytes.Buffer
	countRead  int
	countWrite int
}

func New(buf []byte) *CountingBuffer {
	return &CountingBuffer{
		Buffer: bytes.NewBuffer(buf),
	}
}

func (cb *CountingBuffer) ReadCount() int {
	return cb.countRead
}

func (cb *CountingBuffer) WriteCount() int {
	return cb.countWrite
}

func (cb *CountingBuffer) Read(p []byte) (int, error) {
	n, err := cb.Buffer.Read(p)
	if err != nil {
		return n, err
	}
	cb.countRead += n
	return n, nil
}

func (cb *CountingBuffer) ReadByte() (byte, error) {
	b, err := cb.Buffer.ReadByte()
	if err != nil {
		return b, err
	}
	cb.countRead++
	return b, nil
}

func (cb *CountingBuffer) Write(p []byte) (int, error) {
	n, err := cb.Buffer.Write(p)
	if err != nil {
		return n, err
	}
	cb.countWrite += n
	return n, nil
}

func (cb *CountingBuffer) WriteTo(w io.Writer) (int64, error) {
	n, err := cb.Buffer.WriteTo(w)
	if err != nil {
		return 0, err
	}
	cb.countWrite += int(n)
	return n, nil
}
