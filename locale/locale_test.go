package locale

import (
	"net/http"
	"testing"
)

func TestGet(t *testing.T) {

	t.Run("Пустая локаль", func(t *testing.T) {
		r := &http.Request{
			Header: http.Header{},
		}
		r.Header.Add("Locale-Id", "")
		expected := Ru
		if result := Get(r); result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	t.Run("ru локаль", func(t *testing.T) {
		r := &http.Request{
			Header: http.Header{},
		}
		r.Header.Add("Locale-Id", "ru")
		expected := Ru
		if result := Get(r); result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	t.Run("en локаль", func(t *testing.T) {
		r := &http.Request{
			Header: http.Header{},
		}
		r.Header.Add("Locale-Id", "en")
		expected := En
		if result := Get(r); result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	t.Run("нет хидера с локалью", func(t *testing.T) {
		r := &http.Request{
			Header: http.Header{},
		}
		expected := Ru
		if result := Get(r); result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

}
