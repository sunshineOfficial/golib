package gotime

import (
	"database/sql/driver"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type DateTime struct {
	Date Time
	Time Time
}

func (u DateTime) Merge() time.Time {
	return time.Date(u.Date.Year(), u.Date.Month(), u.Date.Day(), u.Time.Hour(), u.Time.Minute(), u.Time.Second(), u.Time.Nanosecond(), time.UTC)
}

type Time struct {
	time.Time
	OutputLayout string
}

func ToUniversal(t time.Time) Time {
	return Time{
		Time: t,
	}
}

func ToUniversalPointer(t *time.Time) *Time {
	if t == nil {
		return nil
	}

	return &Time{
		Time: *t,
	}
}

// Split конвертирует Time в DateTime (часовой пояс сохраняется)
func (t Time) Split() DateTime {
	year, month, day := t.Date()
	hour, minute, second, nano := t.Hour(), t.Minute(), t.Second(), t.Nanosecond()

	date := ToUniversal(time.Date(year, month, day, 0, 0, 0, 0, t.Location()))
	date.OutputLayout = DateOnly

	time := ToUniversal(time.Date(0, 1, 1, hour, minute, second, nano, t.Location()))
	time.OutputLayout = TimeOnly

	return DateTime{
		Date: date,
		Time: time,
	}
}

func (t *Time) UnmarshalJSON(b []byte) error {
	v := strings.Trim(string(b), "\"")
	return t.parseCleanString(v)
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.format())), nil
}

func (t *Time) UnmarshalXML(d *xml.Decoder, s xml.StartElement) error {
	text := ""
	err := d.DecodeElement(&text, &s)
	if err != nil {
		return err
	}

	return t.parseCleanString(text)
}

func (t Time) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	return e.EncodeElement(t.format(), s)
}

func (t Time) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(t.format())
}

func (t *Time) UnmarshalBSONValue(bsonType bsontype.Type, value []byte) error {
	switch bsonType {
	case bsontype.Null:
		return nil
	case bsontype.String:
		s, _, ok := bsoncore.ReadString(value)
		if !ok {
			return fmt.Errorf("invalid bson string")
		}

		return t.parseCleanString(s)
	case bsontype.DateTime:
		bsonTime, _, ok := bsoncore.ReadTime(value)
		if !ok {
			return fmt.Errorf("invalid bson dateTime")
		}

		t.Time = bsonTime
		return nil
	default:
		return nil
	}
}

func (t *Time) parseCleanString(s string) error {
	if len(s) == 0 || strings.EqualFold(s, "null") {
		return nil
	}

	for _, layout := range _timeLayouts {
		res, err := time.Parse(layout, s)
		if err == nil {
			t.OutputLayout = layout
			t.Time = res
			return nil
		}
	}

	return fmt.Errorf("can't parse time \"%v\" by any known layout", s)
}

func (t Time) format() string {
	if t.IsZero() {
		return ""
	}

	layout := time.RFC3339
	if len(t.OutputLayout) > 0 {
		layout = t.OutputLayout
	}

	return t.Time.Format(layout)
}

func (t *Time) Scan(src any) error {
	if src == nil {
		return nil
	}

	switch v := src.(type) {
	case time.Time:
		t.Time = v
		return nil
	case []byte:
		return t.parseCleanString(string(v))
	case string:
		return t.parseCleanString(v)
	default:
		return fmt.Errorf("cannot scan type %T", src)
	}
}

func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}
