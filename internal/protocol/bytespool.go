package protocol

import (
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
			return make([]byte, 10*1024) // 10KB
		},
	}
	for i := 0; i < 10; i++ {
		pool.Put(pool.New())
	}
	return &BytesPool{
		pool: &pool,
	}
}

func (bp *BytesPool) Get() []byte {
	return bp.pool.Get().([]byte)
}

func (bp *BytesPool) Put(b []byte) {
	bp.pool.Put(b)
}
