package gc

import (
	"bytes"
	"io"

	"github.com/sunshineOfficial/golib/gosync"
)

type BytesBuffer struct {
	buffer *bytes.Buffer
}

func NewBytesBuffer() *BytesBuffer {
	return _bufferPool.Get()
}

func newBuffer(b []byte) *BytesBuffer {
	return &BytesBuffer{
		buffer: bytes.NewBuffer(b),
	}
}

func (b *BytesBuffer) Len() int {
	return b.buffer.Len()
}

func (b *BytesBuffer) Cap() int {
	return b.buffer.Cap()
}

func (b *BytesBuffer) Reset() {
	b.buffer.Reset()
}

func (b *BytesBuffer) Write(p []byte) (int, error) {
	return b.buffer.Write(p)
}

func (b *BytesBuffer) WriteByte(p byte) error {
	return b.buffer.WriteByte(p)
}

func (b *BytesBuffer) WriteString(s string) (int, error) {
	return b.buffer.WriteString(s)
}

func (b *BytesBuffer) String() string {
	return b.buffer.String()
}

func (b *BytesBuffer) Bytes() []byte {
	return b.buffer.Bytes()
}

func (b *BytesBuffer) WriteTo(w io.Writer) (int64, error) {
	return b.buffer.WriteTo(w)
}

func (b *BytesBuffer) ReadFrom(r io.Reader) (int64, error) {
	return b.buffer.ReadFrom(r)
}

func (b *BytesBuffer) Free() {
	if b == nil {
		return
	}

	b.buffer.Reset()
	_bufferPool.Put(b)
}

const (
	MaxPutBufferSize = 16384 // 16 Kb
	InitBufferSize   = 1024  // 1 Kb
)

var _bufferPool = gosync.NewPool[*BytesBuffer](func() *BytesBuffer {
	return newBuffer(make([]byte, 0, InitBufferSize))
}, func(a *BytesBuffer) bool {
	if a == nil {
		return false
	}

	return a.Cap() < MaxPutBufferSize
})
