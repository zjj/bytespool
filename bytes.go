package bytespool

import (
	"container/list"
	"context"
	"errors"
	"io"
	"runtime"
	"sync"
)

/* Bytes */
const (
	unLimitedBytesCapacity int = 0
)

var ErrBytesCapacityTooSmall = errors.New("bytes capacity is too small")

type Bytes struct {
	mux  sync.Mutex
	list *list.List
	pool *Pool

	capacity int
	size     int // current data size that wrote to the list
}

func (bs *Bytes) withPool(p *Pool) *Bytes {
	bs.pool = p
	return bs
}

func (bs *Bytes) capacityInBytes() int {
	return bs.capacity * bs.pool.segmentSize
}

func (bs *Bytes) isUnlimited() bool {
	return bs.capacity == unLimitedBytesCapacity
}

func (bs *Bytes) Write(p []byte) (n int, err error) {
	bs.mux.Lock()
	defer bs.mux.Unlock()

	if !bs.pool.isUnlimited() && bs.size+len(p) > bs.pool.capacityInBytes() {
		return 0, ErrPoolTooSmall
	}

	if !bs.isUnlimited() && bs.size+len(p) > bs.capacityInBytes() {
		return 0, ErrBytesCapacityTooSmall
	}

	if bs.list.Len() == 0 {
		buf, err := bs.pool.Get(context.TODO()) // actually, it wibs never return error
		if err != nil {
			return 0, err
		}
		bs.list.PushBack(buf)
	}

	for len(p) > 0 {
		back := bs.list.Back()
		seg := back.Value.(*segment)

		if seg.size == bs.pool.segmentSize {
			buf, err := bs.pool.Get(context.TODO())
			if err != nil {
				return n, err
			}
			bs.list.PushBack(buf)
			continue
		}

		c := copy(seg.data[seg.size:], p)
		n += c
		bs.size += c
		seg.size += c
		p = p[c:]
	}

	return n, nil
}

func (bs *Bytes) Read(p []byte) (n int, err error) {
	bs.mux.Lock()
	defer bs.mux.Unlock()

	return bs.read(p)
}

func (bs *Bytes) read(p []byte) (n int, err error) {
	for len(p) > 0 && bs.list.Len() > 0 {
		front := bs.list.Front()
		segment := front.Value.(*segment)
		h := min(segment.offset+len(p), segment.size)
		c := copy(p, segment.data[segment.offset:h])
		n += c
		segment.offset += c
		p = p[c:]
		if segment.size == segment.offset {
			bs.list.Remove(front)
			bs.pool.Put(segment)
			// only when a segment is removed from the list,
			// the size should be decreased
			bs.size -= segment.size
		}
	}

	if len(p) > 0 {
		err = io.EOF
	}

	return n, err
}

func (bs *Bytes) ReadAll() ([]byte, error) {
	bs.mux.Lock()
	defer bs.mux.Unlock()

	length := bs.size
	if bs.list.Len() > 0 {
		front := bs.list.Front()
		length = bs.size - front.Value.(*segment).offset
	}
	buf := make([]byte, length)
	_, err := bs.read(buf)
	return buf, err
}

// Free all buffers in the list
// it's not big deal if you forget to call this function
// gc wibs do it for you
func (bs *Bytes) Free() {
	bs.mux.Lock()
	defer bs.mux.Unlock()

	for e := bs.list.Front(); e != nil; e = e.Next() {
		bs.list.Remove(e)
		buf := e.Value.(*segment)
		bs.pool.Put(buf)
	}

	bs.list.Init()
	bs.size = 0
}

func newBytes(capacity int) *Bytes {
	bs := &Bytes{
		list:     list.New(),
		capacity: capacity,
		mux:      sync.Mutex{},
	}

	finalized := func(bs *Bytes) {
		bs.mux.Lock()
		defer bs.mux.Unlock()

		for e := bs.list.Front(); e != nil; e = e.Next() {
			bs.list.Remove(e)
			buf := e.Value.(*segment)
			bs.pool.Put(buf)
		}

		bs.list.Init()
	}

	runtime.SetFinalizer(bs, finalized)
	return bs
}

// BytesPool is a pool of Bytes
type BytesPool struct {
	pool *Pool
}

// NewBytesPool create a new BytesPool
// maxMemory is the threshold value (in byte) of memory the BytesPool could increase,
// eg. 1024 * 1024 * 1024 (1G)
// if maxMemory == 0, the capacity of the pool is unLimited
// segmentSize (byte) is the size of the segment, eg. 1024 * 4,
// maxMeory / segmentSize is the capacity of the pool which means
// the max number of the segments in the pool
func NewBytesPool(maxMemory, segmentSize int) *BytesPool {
	capacity := maxMemory / segmentSize
	if capacity == 0 {
		capacity = unLimitedPoolCapacity
	}
	bp := &BytesPool{
		pool: newPool(capacity, segmentSize),
	}
	return bp
}

// NewBytes create a new Bytes something like []byte
// length is the max length of the Bytes (in byte) eg, 1024 * 1024 (1M)
// data in Bytes is a list of []byte, the size of []byte is pool's segmentSize, like 4k
// so, the max length of the Bytes is bl / segmentSize
// if length == 0, it is not limited, it may use all of the bytes in the pool
func (bp *BytesPool) NewBytes(length int) *Bytes {
	capacity := length / bp.pool.segmentSize
	if length%bp.pool.segmentSize > 0 {
		capacity += 1
	}

	if length == 0 {
		capacity = unLimitedPoolCapacity
	}

	return newBytes(capacity).withPool(bp.pool)
}
