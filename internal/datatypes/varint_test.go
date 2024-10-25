package datatypes

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWriteVarInt(t *testing.T) {
	testCases := []struct {
		inVarInt VarInt
		expOut   []byte
	}{
		{
			inVarInt: 0,
			expOut:   []byte{0x00},
		}, {
			inVarInt: 1,
			expOut:   []byte{0x01},
		}, {
			inVarInt: 128,
			expOut:   []byte{0x80, 0x01},
		}, {
			inVarInt: 255,
			expOut:   []byte{0xff, 0x01},
		}, {
			inVarInt: -1,
			expOut:   []byte{0xff, 0xff, 0xff, 0xff, 0x0f},
		}, {
			inVarInt: -2147483648,
			expOut:   []byte{0x80, 0x80, 0x80, 0x80, 0x08},
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d", tc.inVarInt), func(t *testing.T) {
			out := WriteVarInt(tc.inVarInt)
			require.Equal(t, tc.expOut, out)
		})
	}
}

func TestReadVarInt(t *testing.T) {
	testCases := []struct {
		inBytes []byte
		expOut  VarInt
	}{
		{
			inBytes: []byte{0x00},
			expOut:  0,
		}, {
			inBytes: []byte{0x01},
			expOut:  1,
		}, {
			inBytes: []byte{0x80, 0x01},
			expOut:  128,
		}, {
			inBytes: []byte{0xff, 0x01},
			expOut:  255,
		}, {
			inBytes: []byte{0xff, 0xff, 0xff, 0xff, 0x0f},
			expOut:  -1,
		}, {
			inBytes: []byte{0x80, 0x80, 0x80, 0x80, 0x08},
			expOut:  -2147483648,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%x", tc.inBytes), func(t *testing.T) {
			out, err := ReadVarInt(tc.inBytes)
			require.NoError(t, err)
			require.Equal(t, tc.expOut, out)
		})
	}
}

func TestVarInt(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		in := rand.Int31n(1000)
		t.Run(fmt.Sprintf("%d", in), func(t *testing.T) {
			out, err := ReadVarInt(WriteVarInt(VarInt(in)))
			require.NoError(t, err)
			require.Equal(t, VarInt(in), out)
		})
	}
}

func BenchmarkReadVarInt(b *testing.B) {
	b.ReportAllocs()
	rnd := []byte{0xff, 0x01}
	for i := 0; i < b.N; i++ {
		_, _ = ReadVarInt(rnd)
	}
}

func BenchmarkWriteVarInt(b *testing.B) {
	b.ReportAllocs()
	rand.NewSource(time.Now().Unix())
	rnd := rand.Int31()
	for i := 0; i < b.N; i++ {
		_ = WriteVarInt(VarInt(rnd))
	}
}
