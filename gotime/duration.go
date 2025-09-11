package gotime

import (
	"strings"
	"time"
)

type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	text := strings.Trim(string(b), ` '"`)
	duration, err := time.ParseDuration(text)
	if err != nil {
		return err

	}

	*d = Duration(duration)
	return nil
}
