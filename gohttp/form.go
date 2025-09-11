package gohttp

import (
	"mime/multipart"
	"net/textproto"
)

type FormDataField struct {
	Name  string
	Value string
}

type FormDataFile struct {
	Payload    multipart.File
	MIMEHeader textproto.MIMEHeader
}
