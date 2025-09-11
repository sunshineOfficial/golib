package gotime

import "time"

// Days расчитывает разницу дней между датами
// не учитывает время,
// например если передать параметры
// d1 = 15-01-2018 15:00
// d2 = 16-01-2018 00:00
// то ответ будет 1 день
func Days(d1, d2 time.Time) float64 {
	loc := time.Now().Location()
	v1 := time.Date(d1.Year(), d1.Month(), d1.Day(), 0, 0, 0, 0, loc)
	v2 := time.Date(d2.Year(), d2.Month(), d2.Day(), 0, 0, 0, 0, loc)
	return float64(int(v1.Sub(v2).Hours() / 24))
}
