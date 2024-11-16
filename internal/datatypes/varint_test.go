package datatypes

import (
	"bytes"
	"fmt"
	"github.com/BinaryArchaism/mc-srv/internal/countingbuffer"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestVarInt_Write(t *testing.T) {
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
			buf := countingbuffer.New(make([]byte, 0, 10))
			err := tc.inVarInt.Write(buf)
			require.NoError(t, err)
			require.Equal(t, tc.expOut, buf.Bytes())
		})
	}
}

func TestVarInt_Read(t *testing.T) {
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
			var varInt VarInt
			err := varInt.Read(bytes.NewReader(tc.inBytes))
			require.NoError(t, err)
			require.Equal(t, tc.expOut, varInt)
		})
	}
}

func TestVarInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		rnd := rand.Int31()
		varInt := VarInt(rnd)
		buf := countingbuffer.New(make([]byte, 0, 5))
		err := varInt.Write(buf)
		require.NoError(t, err, i)

		err = varInt.Read(buf)
		require.NoError(t, err, i)
		require.Equal(t, rnd, int32(varInt), i)
	}
}

func BenchmarkReadVarInt(b *testing.B) {
	b.ReportAllocs()
	rnd := []byte{0xff, 0x01}
	var varInt VarInt
	for i := 0; i < b.N; i++ {
		_ = varInt.Read(bytes.NewBuffer(rnd))
	}
}

func BenchmarkWriteVarInt(b *testing.B) {
	b.ReportAllocs()
	varInt := VarInt(12345)
	for i := 0; i < b.N; i++ {
		_ = varInt.Write(bytes.NewBuffer([]byte{}))
	}
}

//
//func TestWriteVarInt(t *testing.T) {
//	testCases := []struct {
//		inVarInt VarInt
//		expOut   []byte
//	}{
//		{
//			inVarInt: 0,
//			expOut:   []byte{0x00},
//		}, {
//			inVarInt: 1,
//			expOut:   []byte{0x01},
//		}, {
//			inVarInt: 128,
//			expOut:   []byte{0x80, 0x01},
//		}, {
//			inVarInt: 255,
//			expOut:   []byte{0xff, 0x01},
//		}, {
//			inVarInt: -1,
//			expOut:   []byte{0xff, 0xff, 0xff, 0xff, 0x0f},
//		}, {
//			inVarInt: -2147483648,
//			expOut:   []byte{0x80, 0x80, 0x80, 0x80, 0x08},
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(fmt.Sprintf("%d", tc.inVarInt), func(t *testing.T) {
//			out := WriteVarInt(tc.inVarInt)
//			require.Equal(t, tc.expOut, out)
//		})
//	}
//}
//func TestBinaryWriteVarInt(t *testing.T) {
//	testCases := []struct {
//		inVarInt VarInt
//		expOut   []byte
//	}{
//		{
//			inVarInt: 0,
//			expOut:   []byte{0x00},
//		}, {
//			inVarInt: 1,
//			expOut:   []byte{0x01},
//		}, {
//			inVarInt: 128,
//			expOut:   []byte{0x80, 0x01},
//		}, {
//			inVarInt: 255,
//			expOut:   []byte{0xff, 0x01},
//		}, {
//			inVarInt: -1,
//			expOut:   []byte{0xff, 0xff, 0xff, 0xff, 0x0f},
//		}, {
//			inVarInt: -2147483648,
//			expOut:   []byte{0x80, 0x80, 0x80, 0x80, 0x08},
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(fmt.Sprintf("%d", tc.inVarInt), func(t *testing.T) {
//			out := BinaryWriteVarInt(int(tc.inVarInt))
//			require.Equal(t, tc.expOut, out)
//		})
//	}
//}
//
//func TestReadVarInt(t *testing.T) {
//	testCases := []struct {
//		inBytes []byte
//		expOut  VarInt
//	}{
//		{
//			inBytes: []byte{0x00},
//			expOut:  0,
//		}, {
//			inBytes: []byte{0x01},
//			expOut:  1,
//		}, {
//			inBytes: []byte{0x80, 0x01},
//			expOut:  128,
//		}, {
//			inBytes: []byte{0xff, 0x01},
//			expOut:  255,
//		}, {
//			inBytes: []byte{0xff, 0xff, 0xff, 0xff, 0x0f},
//			expOut:  -1,
//		}, {
//			inBytes: []byte{0x80, 0x80, 0x80, 0x80, 0x08},
//			expOut:  -2147483648,
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(fmt.Sprintf("%x", tc.inBytes), func(t *testing.T) {
//			out, err := ReadVarInt(tc.inBytes)
//			require.NoError(t, err)
//			require.Equal(t, tc.expOut, out)
//		})
//	}
//}
//func TestBinaryReadVarInt(t *testing.T) {
//	testCases := []struct {
//		inBytes []byte
//		expOut  VarInt
//	}{
//		{
//			inBytes: []byte{0x00},
//			expOut:  0,
//		}, {
//			inBytes: []byte{0x01},
//			expOut:  1,
//		}, {
//			inBytes: []byte{0x80, 0x01},
//			expOut:  128,
//		}, {
//			inBytes: []byte{0xff, 0x01},
//			expOut:  255,
//		}, {
//			inBytes: []byte{0xff, 0xff, 0xff, 0xff, 0x0f},
//			expOut:  -1,
//		}, {
//			inBytes: []byte{0x80, 0x80, 0x80, 0x80, 0x08},
//			expOut:  -2147483648,
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(fmt.Sprintf("%x", tc.inBytes), func(t *testing.T) {
//			buf := bytes.NewBuffer(tc.inBytes)
//			out, err := BinaryReadVarInt(buf)
//			require.NoError(t, err)
//			require.Equal(t, tc.expOut, out)
//		})
//	}
//}
//
//func TestVarInt(t *testing.T) {
//	for i := 0; i < 10_000; i++ {
//		in := rand.Int31n(1000)
//		t.Run(fmt.Sprintf("%d", in), func(t *testing.T) {
//			out, err := ReadVarInt(WriteVarInt(VarInt(in)))
//			require.NoError(t, err)
//			require.Equal(t, VarInt(in), out)
//		})
//	}
//}
//
//func BenchmarkReadVarInt(b *testing.B) {
//	b.ReportAllocs()
//	rnd := []byte{0xff, 0x01}
//	for i := 0; i < b.N; i++ {
//		_, _ = ReadVarInt(rnd)
//	}
//}
//
//func BenchmarkBinaryReadVarInt(b *testing.B) {
//	b.ReportAllocs()
//	buf := countingbuffer.New([]byte{0xff, 0x01})
//	for i := 0; i < b.N; i++ {
//		_, _ = BinaryReadVarInt(buf)
//	}
//}
//
//func BenchmarkWriteVarInt(b *testing.B) {
//	b.ReportAllocs()
//	rand.NewSource(time.Now().Unix())
//	rnd := rand.Int31()
//	for i := 0; i < b.N; i++ {
//		_ = WriteVarInt(VarInt(rnd))
//	}
//}
//
//func BenchmarkBinaryWriteVarInt(b *testing.B) {
//	b.ReportAllocs()
//	rand.NewSource(time.Now().Unix())
//	rnd := rand.Int31()
//	for i := 0; i < b.N; i++ {
//		_ = BinaryWriteVarInt(int(rnd))
//	}
//}
