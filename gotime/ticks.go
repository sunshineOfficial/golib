package gotime

import (
	"time"
)

const (
	_ticksOffset   = 621355968000000000
	_ticksInSecond = 10000000
)

// TimeToTicks преобразовывает дату из Go в тики .Net (Windows)
func TimeToTicks(t time.Time) int64 {
	return (t.Unix() * _ticksInSecond) + _ticksOffset
}

// TicksToTime преобразовывает тики .Net (Windows) в дату из Go
func TicksToTime(ticks int64) time.Time {
	return time.Unix((ticks-_ticksOffset)/_ticksInSecond, 0)
}

// ParseDateTimeTicks получает тики из строки вида "02.01.2006 15:04"
func ParseDateTimeTicks(value string) (int64, error) {
	t, err := time.Parse(DateTimeNet, value)
	return TimeToTicks(t), err
}

// ParseDateTicks получает тики из строки вида "02.01.2006"
func ParseDateTicks(value string) (int64, error) {
	t, err := time.Parse(DateOnlyNet, value)
	return TimeToTicks(t), err
}

// ParseTimeTicks получает тики из строки вида "15:04"
func ParseTimeTicks(value string) (int64, error) {
	t, err := time.Parse(TimeOnlyNet, value)
	return TimeToTicks(t), err
}
