package gotime

import (
	"log"
	"time"
)

// MustConvertToNet преобразовывает время из формата yyyy-mm-dd в dd.mm.yyyy
// При несоответствии входной строки - возвращает пустую строку
func MustConvertToNet(date string) string {
	t, err := time.Parse(DateOnly, date)

	if err != nil {
		log.Println(err)
		return ""
	}

	return t.Format(DateOnlyNet)
}
