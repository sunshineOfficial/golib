package language

import (
	"testing"
)

func TestKeyboardRuToEn(t *testing.T) {
	result := SwitchKeys("привет")
	if result != "ghbdtn" {
		t.Error("Error switch \"привет\" to \"ghbdtn\"")
	}
}

func TestKeyboardEnToRu(t *testing.T) {
	result := SwitchKeys("ghbdtn")
	if result != "привет" {
		t.Error("Error switch \"ghbdtn\" to \"привет\"")
	}
}

func TestKeyboardSpecSymbols(t *testing.T) {
	result := SwitchKeys("хъжэбю")
	if result != "[];',." {
		t.Error("Error switch \"хъжэбю\" to \"[];',.\"")
	}
}

func TestKeyboardShiftSpecSymbols(t *testing.T) {
	result := SwitchKeys("ХЪЖЭБЮ")
	if result != "{}:\"<>" {
		t.Error("Error switch \"хъжэбю\" to \"{}:\"<>\"")
	}
}
