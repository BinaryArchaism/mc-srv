package datatypes

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestShort(t *testing.T) {
	for i := 0; i < 1000; i++ {
		rnd := int16(rand.Int31())
		res := ReadShort(WriteShort(Short(rnd)))
		require.Equal(t, rnd, int16(res))
	}
}

func TestUShort(t *testing.T) {
	for i := 0; i < 1000; i++ {
		rnd := uint16(rand.Int31())
		res := ReadUShort(WriteUShort(UShort(rnd)))
		require.Equal(t, rnd, uint16(res))
	}
}
