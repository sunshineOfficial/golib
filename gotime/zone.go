package gotime

import (
	"fmt"
	"time"
)

var (
	Moscow, _ = time.LoadLocation("Europe/Moscow")
)

func MoscowNow() time.Time {
	return time.Now().In(Moscow)
}

func AtLocation(t time.Time, location *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, location)
}

func AtLocationByName(t time.Time, name string) (time.Time, error) {
	location, err := time.LoadLocation(name)
	if err != nil {
		return time.Time{}, fmt.Errorf("ошибка загрузки часового пояса %s: %s", name, err)
	}

	return AtLocation(t, location), nil
}

func IsSameLocation(t1, t2 time.Time) bool {
	return t1.Location().String() == t2.Location().String()
}

// GetUtcOffset возвращает разницу между местным временем и UTC
func GetUtcOffset(loc *time.Location) time.Duration {
	_, offset := time.Now().In(loc).Zone()
	return time.Second * time.Duration(offset)
}

// NewUnixTime возвращает экземпляр time.Time содержащий дату полученную из unix в поясе loc
// стандартный метод time.Unix возвращает время в местном часовом поясе
func NewUnixTime(unix int64, loc *time.Location) time.Time {
	return ReplaceLocation(time.Unix(unix, 0).Add(GetUtcOffset(time.Local)*-1), loc)
}

// ReplaceLocation подменяет часовой пояс во времени t, значения даты и времени не изменяются
func ReplaceLocation(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
}
