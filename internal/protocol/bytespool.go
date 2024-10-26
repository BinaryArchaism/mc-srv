package protocol

import (
	"bytes"
	"sync"
)

// todo pool for different sizes

var Pool = NewBytesPool()

type BytesPool struct {
	pool *sync.Pool
}

func NewBytesPool() *BytesPool {
	pool := sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 10*1024)) // 10KB
		},
	}
	for i := 0; i < 10; i++ {
		pool.Put(pool.New())
	}
	return &BytesPool{
		pool: &pool,
	}
}

func (bp *BytesPool) Get() *bytes.Buffer {
	return bp.pool.Get().(*bytes.Buffer)
}

func (bp *BytesPool) Put(b *bytes.Buffer) {
	b.Reset()
	bp.pool.Put(b)
}
