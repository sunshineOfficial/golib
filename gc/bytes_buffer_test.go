package gc

import (
	"testing"
)

func Test(t *testing.T) {
	src := NewBytesBuffer()
	defer src.Free()

	src.WriteString("hello")
	src.WriteByte(',')
	src.Write([]byte(" world"))

	got := src.String()
	want := "hello, world"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAlloc(t *testing.T) {
	got := int(testing.AllocsPerRun(5, func() {
		src := NewBytesBuffer()
		defer src.Free()

		src.WriteString("not 1K worth of bytes")
	}))
	if got != 0 {
		t.Errorf("got %d allocs, want 0", got)
	}
}
