package goio

import (
	"bytes"
	"io"
)

func NewBytesReader(b []byte) io.Reader {
	n := make([]byte, len(b))
	copy(n, b)
	return bytes.NewReader(n)
}

func NewBytesReadCloser(b []byte) io.ReadCloser {
	return io.NopCloser(NewBytesReader(b))
}

func NewReadCloser(r io.Reader) io.ReadCloser {
	return io.NopCloser(r)
}

func ReadString(r io.Reader) (string, error) {
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(r)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
