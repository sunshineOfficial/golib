package language

import (
	"regexp"
)

func TranslitEnToRu(text string) string {
	cyrillicString := text
	for _, pair := range enTranslations {
		r, _ := regexp.Compile(pair.K)
		cyrillicString = r.ReplaceAllString(cyrillicString, pair.V)
	}
	return cyrillicString
}

func TranslitRuToEn(text string) string {
	letters := []rune(text)
	var result []rune
	for _, letter := range letters {
		if translit, ok := ruTranslations[letter]; ok {
			result = append(result, []rune(translit)...)
			continue
		}
		result = append(result, letter)
	}
	return string(result)
}
