package language

import "unicode"

const (
	RU   = Lang(1 << 0)
	EN   = Lang(1 << 1)
	RUEN = RU | EN
	ENRU = RU | EN
)

type Lang int

func IdentifyOne(letter rune) (lang Lang) {
	if unicode.Is(unicode.Cyrillic, letter) {
		lang = RU
	}

	if unicode.Is(unicode.Latin, letter) {
		lang = EN
	}
	return
}

func Identify(str string) (lang Lang) {
	letters := []rune(str)
	for _, letter := range letters {
		lang |= IdentifyOne(letter)
	}
	return
}

func IsOneRU(letter rune) bool {
	return RU == IdentifyOne(letter)
}

func IsOneEN(letter rune) bool {
	return EN == IdentifyOne(letter)
}

func IsRU(str string) bool {
	return RU == Identify(str)
}

func IsEN(str string) bool {
	return EN == Identify(str)
}
