package datatypes

import (
	"bytes"
	"compress/zlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestString(t *testing.T) {
	for i := 0; i < 100; i++ {
		rndStr := randStringRunes(rand.Intn(100))
		ctmStr := FromString(rndStr)
		bytesStr := WriteString(ctmStr)
		ctmStr = ReadString(bytesStr)
		outStr := ToString(ctmStr)
		assert.Equal(t, rndStr, outStr)
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

//FE01FA000B004D0043007C00500069006E00670048006F0073007400197F0009006C006F00630061006C0068006F0073007400001F90

func TestReadString(t *testing.T) {
	bstr := "FA000B004D0043007C00500069006E00670048006F0073007400197F0009006C006F00630061006C0068006F0073007400001F90"
	var b []byte
	for i := 0; i < len(bstr); i++ {

	}

	uncompressed, err := zlib.NewReader(bytes.NewReader(b))
	require.NoError(t, err)

	t.Log(uncompressed)
}
