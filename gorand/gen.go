package gorand

import (
	"crypto/rand"
	"math/big"
)

func RandomInt(min, max int) int {
	value, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return int(value.Int64()) + min
}

// RandomString возвращает случайную строку длиной от min до max символов. Для создания строки фиксированной длины
// нужно указать одинаковые min и max
func RandomString(alphabet []rune, min, max int) string {
	if alphabet == nil {
		alphabet = DefaultAlphabet
	}

	count := min
	if min != max {
		count = RandomInt(min, max)
	}

	result := make([]rune, count)
	for i := 0; i < count; i++ {
		runeIdx := RandomInt(0, len(alphabet))
		result[i] = alphabet[runeIdx]
	}

	return string(result)
}
