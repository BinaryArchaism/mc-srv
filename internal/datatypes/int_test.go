package datatypes

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestUint16(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		rnd := uint16(rand.Int31())
		res := ReadUint16(WriteUint16(Uint16(rnd)))
		require.Equal(t, rnd, uint16(res))
	}
}
