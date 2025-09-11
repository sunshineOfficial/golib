package gohttp

const (
	ContentTypeHeader     = "Content-Type"
	ContentLengthHeader   = "Content-Length"
	UserAgentHeader       = "User-Agent"
	ContentEncodingHeader = "Content-Encoding"
	AcceptEncodingHeader  = "Accept-Encoding"
	AuthorizationHeader   = "Authorization"
	OriginHeader          = "Origin"
	RequestIdHeader       = "X-Request-Id"
)

const (
	ContentEncodingEmpty = ""
	ContentEncodingGzip  = "gzip"
)

const (
	ContentTypeEmpty = ""

	ContentTypeJSON = "application/json"
	ContentTypeXML  = "application/xml"
	ContentTypeHTML = "text/html"
	ContentTypeText = "text/plain"
)

const (
	UserAgentDefault = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:104.0) Gecko/20100101 Firefox/104.0"
)
