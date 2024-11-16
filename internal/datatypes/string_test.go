package datatypes

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestString_Read(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x03, 'a', 'b', 'c'})
	var s String
	err := s.Read(buf)
	require.NoError(t, err)
	require.Equal(t, "abc", s.String())
}

func TestString_Write(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	s := String{
		Size: 3,
		Data: "abc",
	}
	err := s.Write(buf)
	require.NoError(t, err)
	require.Equal(t, []byte{0x03, 'a', 'b', 'c'}, buf.Bytes())
}

func TestString(t *testing.T) {
	for i := 0; i < 100; i++ {
		rndStr := randStringRunes(rand.Intn(100))
		var s String
		s.FromString(rndStr)
		require.Equal(t, rndStr, s.String())

		buf := bytes.NewBuffer([]byte{})
		err := s.Write(buf)
		require.NoError(t, err)

		var s2 String
		buf = bytes.NewBuffer(buf.Bytes())
		err = s2.Read(buf)
		require.NoError(t, err)

		require.Equal(t, s, s2)
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
