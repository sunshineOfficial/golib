package gotime

import "time"

const (
	RFC3339WithoutLocation = "2006-01-02T15:04:05"
	DateTimeWithSpace      = "2006-01-02 15:04:05"

	TimeOnly = "15:04:05"
	DateOnly = "2006-01-02"

	TimeOnlyNet = "15:04"
	DateOnlyNet = "02.01.2006"
	DateTimeNet = "02.01.2006 15:04"
)

var (
	_timeLayouts = []string{
		DateTimeWithSpace, RFC3339WithoutLocation, time.RFC3339Nano, time.RFC3339,
		time.RFC822, time.RFC822Z, time.RFC1123, time.RFC1123Z, time.RFC850,
		TimeOnly, DateOnly,
	}
)
