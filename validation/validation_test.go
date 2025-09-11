package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsRequireChars(t *testing.T) {
	t.Run("string contains upper letters", func(t *testing.T) {
		expected := true
		value := HasUpper("TSugar28!")
		assert.Equal(t, expected, value)
	})
	t.Run("string doesnt contains upper letters", func(t *testing.T) {
		expected := false
		value := HasUpper("redonu")
		assert.Equal(t, expected, value)
	})
	t.Run("string contains lower letters", func(t *testing.T) {
		expected := true
		value := HasLower("TSugar28!")
		assert.Equal(t, expected, value)
	})
	t.Run("string doesnt contains lower letters", func(t *testing.T) {
		expected := false
		value := HasLower("REDONU")
		assert.Equal(t, expected, value)
	})
	t.Run("string contains numbers", func(t *testing.T) {
		expected := true
		value := HasNumeric("TSugar28!")
		assert.Equal(t, expected, value)
	})
	t.Run("string doesnt contains numbers", func(t *testing.T) {
		expected := false
		value := HasNumeric("redonu")
		assert.Equal(t, expected, value)
	})
	t.Run("string contains special chars", func(t *testing.T) {
		expected := true
		value := HasSpecial("TSugar28!")
		assert.Equal(t, expected, value)
	})
	t.Run("string doesnt contains special chars", func(t *testing.T) {
		expected := false
		value := HasSpecial("redonu")
		assert.Equal(t, expected, value)
	})
}

func TestLengthLess(t *testing.T) {
	t.Run("string length is correct", func(t *testing.T) {
		expected := false
		value := LengthLess("redonu28", 8)
		assert.Equal(t, expected, value)
	})
	t.Run("string length is incorrect", func(t *testing.T) {
		expected := true
		value := LengthLess("1234", 8)
		assert.Equal(t, expected, value)
	})
}

func TestContainsUsernamePart(t *testing.T) {
	t.Run("string doesnt contains username", func(t *testing.T) {
		expected := false
		value := ContainsPart("v.stepanchev@sunshine.today", "TSugar28!", 2)
		assert.Equal(t, expected, value)
	})
	t.Run("string contains username", func(t *testing.T) {
		expected := true
		value := ContainsPart("v.stepanchev@sunshine.today", "ste", 2)
		assert.Equal(t, expected, value)
	})
	t.Run("string contains username", func(t *testing.T) {
		expected := true
		value := ContainsPart("v.stepanchev@sunshine.today", "v.s", 2)
		assert.Equal(t, expected, value)
	})
	t.Run("string contains username", func(t *testing.T) {
		expected := true
		value := ContainsPart("v.stepanchev@sunshine.today", "ev@", 2)
		assert.Equal(t, expected, value)
	})
}
