package language

func SwitchKeys(str string) string {
	letters := []rune(str)
	for i, letter := range letters {
		if key, ok := keys[letter]; ok {
			letters[i] = key
		}
	}

	return string(letters)
}
