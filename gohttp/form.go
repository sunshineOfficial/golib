package gohttp

import (
	"io"
)

// MultipartField представляет поле формы
type MultipartField struct {
	Name  string
	Value string
}

// MultipartFile представляет файл для загрузки
type MultipartFile struct {
	FieldName string
	FileName  string
	Reader    io.Reader // Содержимое файла
}

// MultipartData содержит данные для multipart формы
type MultipartData struct {
	Fields []MultipartField
	Files  []MultipartFile
}
