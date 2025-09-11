package locale

import "net/http"

const (
	Header = "Locale-Id"
)

var (
	Ru = "ru"
	En = "en"
)

func Get(r *http.Request) string {
	return GetFromHeader(r.Header)
}

func GetFromHeader(header http.Header) string {
	l := header.Get(Header)
	if l == "" {
		return Ru
	}
	return l
}

func Set(r *http.Request, locale string) {
	SetToHeader(r.Header, locale)
}

func SetToHeader(header http.Header, locale string) {
	header.Set(Header, locale)
}
