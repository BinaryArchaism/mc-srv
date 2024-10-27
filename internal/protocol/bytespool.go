package protocol

import (
	"sync"
	"sync/atomic"
)

const (
	poolBucketCount = 3

	SmallObjectSize  = 1024        // 1KB
	MediumObjectSize = 10 * 1024   // 10KB
	LargeObjectSize  = 1024 * 1024 // 1MB
)

var Pool = NewBytesPool()

type BytesPool struct {
	buckets   [poolBucketCount]*bucket
	allocated atomic.Int32
}

func NewBytesPool() *BytesPool {
	return &BytesPool{
		buckets: [poolBucketCount]*bucket{
			newBucket(SmallObjectSize, 1000),
			newBucket(MediumObjectSize, 100),
			newBucket(LargeObjectSize, 10),
		},
	}
}

// Stats return map[poolSize]inUseObjects, poolSize = -1 shows all allocation outside of pools
func (bp *BytesPool) Stats() map[int]int32 {
	m := make(map[int]int32, poolBucketCount+1)
	for _, v := range bp.buckets {
		m[v.size] = v.inUse.Load()
	}
	m[-1] = bp.allocated.Load()
	return m
}

// GetN return pre allocated bytes from pool if size < size of objPool elem
// or allocate new by make([]byte, size)
func (bp *BytesPool) GetN(size int) []byte {
	for _, v := range bp.buckets {
		if size <= v.size {
			return v.Get()
		}
	}
	bp.allocated.Add(1)
	return make([]byte, size)
}

// Put stores b in pre allocated pool of byte or ignores if obj too large
func (bp *BytesPool) Put(b []byte) {
	for _, v := range bp.buckets {
		if cap(b) <= v.size {
			v.Put(b)
		}
	}
}

type bucket struct {
	inUse atomic.Int32
	size  int
	pool  *sync.Pool
}

func newBucket(size int, count int) *bucket {
	p := bucket{
		size: size,
		pool: &sync.Pool{
			New: func() any {
				return make([]byte, size)
			},
		},
	}
	for i := 0; i < count; i++ {
		p.pool.Put(p.pool.New())
	}
	return &p
}

func (pool *bucket) Get() []byte {
	pool.inUse.Add(1)
	return pool.pool.Get().([]byte)
}

func (pool *bucket) Put(b []byte) {
	pool.inUse.Add(-1)
	pool.pool.Put(b)
}
