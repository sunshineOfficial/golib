package goio

import (
	"compress/gzip"
	"io"
	"net/http"
)

type Compression string

const (
	CompressionEmpty Compression = ""
	CompressionGzip  Compression = "gzip"

	ContentEncodingHeader = "Content-Encoding"
)

// GetCompression возвращает сжатие Compression соответствущее значению заголовка ContentEncodingHeader из header
func GetCompression(header http.Header) Compression {
	contentEncoding := header.Get(ContentEncodingHeader)
	if len(contentEncoding) > 0 {
		return Compression(contentEncoding)
	}

	return CompressionEmpty
}

func Decompress(compression Compression, origin io.Reader) (io.Reader, error) {
	switch compression {
	case CompressionGzip:
		return gzip.NewReader(origin)
	default:
		return origin, nil
	}
}

func Compress(compression Compression, writer io.Writer) io.Writer {
	switch compression {
	case CompressionGzip:
		return gzip.NewWriter(writer)
	default:
		return writer
	}
}
