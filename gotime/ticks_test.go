package gotime

import (
	"testing"
)

func TestGetTicksFromDate(t *testing.T) {
	tick, err := ParseDateTicks("01.01.2018")
	if err != nil {
		t.Error(err)
	}
	if tick != 636503616000000000 {
		t.Error("tick not equals", tick)
	}
}

func TestTicksToTime(t *testing.T) {
	date := TicksToTime(636503616000000000)
	if date.String() != "2018-01-01 03:00:00 +0300 MSK" {
		t.Error("не правильно считает перевод тиков в time.Time ", date.String())
	}
}
