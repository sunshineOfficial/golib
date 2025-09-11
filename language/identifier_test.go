package language

import (
	"testing"
)

func TestIdentifyOneEN(t *testing.T) {
	lang := IdentifyOne('b')
	if lang != EN {
		t.Error("Error identify EN-letter")
	}
}

func TestIdentifyOneRU(t *testing.T) {
	lang := IdentifyOne('б')
	if lang != RU {
		t.Error("Error identify RU-letter")
	}
}

func TestIdentifySomeEN(t *testing.T) {
	lang := Identify("some text")
	if lang != EN {
		t.Error("Error identify EN text")
	}
}

func TestIdentifySomeRU(t *testing.T) {
	lang := Identify("некоторый текст")
	if lang != RU {
		t.Error("Error identify RU text")
	}
}

func TestIsOneRU(t *testing.T) {
	if !IsOneRU('б') {
		t.Error("IsOneRU: Error identify ru letter")
	}
	if IsOneRU('b') {
		t.Error("IsOneRU: Error identify en letter")
	}
}

func TestIsOneEN(t *testing.T) {
	if !IsOneEN('b') {
		t.Error("IsOneEN: Error identify en letter")
	}
	if IsOneEN('б') {
		t.Error("IsOneEN: Error identify ru letter")
	}
}

func TestIsAllRU(t *testing.T) {
	if !IsRU("некоторый текст") {
		t.Error("IsRU: Error identify ru letters")
	}
	if IsRU("some text") {
		t.Error("IsRU: Error identify en letters")
	}
}

func TestIsAllEN(t *testing.T) {
	if !IsEN("some text") {
		t.Error("IsEN: Error identify en letters")
	}
	if IsEN("некоторый текст") {
		t.Error("IsEN: Error identify ru letters")
	}

}
