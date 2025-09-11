package validation

import (
	"regexp"
	"strings"
)

func HasNumeric(str string) bool {
	containsNumeric := regexp.MustCompile(`\d`).MatchString
	return containsNumeric(str)
}

func HasUpper(str string) bool {
	containsUpper := regexp.MustCompile(`[A-Z]`).MatchString
	return containsUpper(str)
}

func HasLower(str string) bool {
	containsLower := regexp.MustCompile(`[a-z]`).MatchString
	return containsLower(str)
}

func HasSpecial(str string) bool {
	containsSpecial := regexp.MustCompile(`[!@#%\^&*()_+={}|;':",.\-\[\]<>?~\\]`).MatchString
	return containsSpecial(str)
}

func LengthLess(str string, length int) bool {
	return len(str) < length
}

func ContainsPart(source, part string, partLength int) bool {
	for i := 0; i < len(source)-partLength; i++ {
		te := source[i : i+partLength+1]
		if strings.Contains(part, te) {
			return true
		}
	}

	return false
}
