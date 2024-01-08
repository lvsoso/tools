package storage

import (
	"io"
	"sync/atomic"
)

type ReaderCloserCounter struct {
	count uint64
	io.Reader
}

func NewReaderCloserCounter(r io.Reader) *ReaderCloserCounter {
	return &ReaderCloserCounter{
		Reader: r,
	}
}

func (counter *ReaderCloserCounter) Read(buf []byte) (int, error) {
	n, err := counter.Reader.Read(buf)
	if n >= 0 {
		atomic.AddUint64(&counter.count, uint64(n))
	}

	return n, err
}

func (counter *ReaderCloserCounter) Count() uint64 {
	return atomic.LoadUint64(&counter.count)
}

func (counter *ReaderCloserCounter) Close() error {
	return nil
}
