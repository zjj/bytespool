package bytespool

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

const (
	unLimitedPoolCapacity int = 0
	defaultSegmentSize    int = 1024 * 4
)

type segment struct {
	data   []byte
	size   int // size of data with content in the segment, size would always increase
	offset int // offset of data in the segment, it is used for read
}

var ErrPoolFull = errors.New("pool is full")
var ErrPoolTooSmall = errors.New("pool is too small")

type Pool struct {
	sem           *semaphore.Weighted
	pool          sync.Pool
	capacity      int
	curSegmentNum int64
	segmentSize   int
}

func (bp *Pool) capacityInBytes() int {
	return bp.capacity * bp.segmentSize
}

func (bp *Pool) Get(ctx context.Context) (*segment, error) {
	if bp.capacity != unLimitedPoolCapacity {
		if err := bp.sem.Acquire(ctx, 1); err != nil {
			return nil, err
		}
	}

	atomic.AddInt64(&bp.curSegmentNum, 1)
	return bp.pool.Get().(*segment), nil
}

func (bp *Pool) Put(seg *segment) {
	if bp.capacity != unLimitedPoolCapacity {
		bp.sem.Release(1)
	}

	// check if the buf is from Pool, do not put it back to pool
	if len(seg.data) != bp.segmentSize {
		return
	}

	atomic.AddInt64(&bp.curSegmentNum, -1)
	seg.size = 0
	seg.offset = 0
	bp.pool.Put(seg)
}

// newPool creates a new buffer pool.
// capacity is the number of buffer in the pool.
// segSize is the size of each buffer.
// If capacity is 0, the pool is unlimit.
func newPool(capacity int, segSize int) *Pool {
	if capacity < 0 {
		panic("capacity must be greater than or equal to 0")
	}

	if segSize <= 0 {
		panic("segSize must be greater than 0")
	}

	bp := &Pool{
		pool: sync.Pool{
			New: func() interface{} {
				return &segment{
					data: make([]byte, segSize),
					size: 0,
				}
			},
		},
		segmentSize: segSize,
		capacity:    capacity,
	}

	if capacity != unLimitedPoolCapacity {
		bp.sem = semaphore.NewWeighted(int64(capacity))
	}

	return bp
}
