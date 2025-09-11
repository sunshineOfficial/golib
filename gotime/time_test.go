package gotime

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

// nolint
type testWrapper struct {
	XMLName xml.Name `json:"-" xml:"Wrapper"`
	Time    Time     `json:"time" xml:"time"`
	TimeNil *Time    `json:"timeNil,omitempty" xml:"timeNil,omitempty"`
}

type testCase struct {
	Time   time.Time
	Layout string
	Text   string
}

var (
	msk, _ = time.LoadLocation("Europe/Moscow")
	gmt, _ = time.LoadLocation("Etc/GMT")

	testCases = []testCase{
		{
			Layout: time.RFC3339,
			Text:   "2022-06-20T08:55:11Z",
			Time:   time.Date(2022, 6, 20, 8, 55, 11, 0, time.UTC),
		},
		{
			Layout: "2006-01-02T15:04:05",
			Text:   "2022-06-20T08:55:11",
			Time:   time.Date(2022, 6, 20, 8, 55, 11, 0, time.Local),
		},
		{
			Layout: "2006-01-02 15:04:05",
			Text:   "2022-06-20 08:55:11",
			Time:   time.Date(2022, 6, 20, 8, 55, 11, 0, time.Local),
		},
		{
			Layout: time.RFC822,
			Text:   "29 Jun 22 08:55 GMT",
			Time:   time.Date(2022, 6, 29, 8, 55, 0, 0, gmt),
		},
		{
			Layout: time.RFC822Z,
			Text:   "29 Jun 22 08:55 +0300",
			Time:   time.Date(2022, 6, 29, 8, 55, 0, 0, msk),
		},
		{
			Layout: TimeOnly,
			Text:   "10:13:00",
			Time:   time.Date(0, 1, 1, 10, 13, 0, 0, time.Local),
		},
	}
)

func TestUniversalTime_UnmarshalJson(t *testing.T) {
	for _, testCase := range testCases {
		w := testWrapper{}
		assert.Nil(t, json.Unmarshal([]byte(createJson(testCase.Text)), &w))
		testTime(t, testCase.Text, testCase.Time, w.Time.Time)
	}
}

func TestUniversalTime_MarshalJson(t *testing.T) {
	for _, testCase := range testCases {
		w := testWrapper{
			Time: Time{
				Time:         testCase.Time,
				OutputLayout: testCase.Layout,
			},
		}
		b, err := json.Marshal(w)
		assert.Nil(t, err)

		assert.Equal(t, createJson(testCase.Text), string(b), testCase.Layout)
	}
}

func TestUniversalTime_UnmarshalXml(t *testing.T) {
	for _, testCase := range testCases {
		w := testWrapper{}
		assert.Nil(t, xml.Unmarshal([]byte(createXml(testCase.Text)), &w))
		testTime(t, testCase.Text, testCase.Time, w.Time.Time)
	}
}

func TestUniversalTime_MarshalXml(t *testing.T) {
	for _, testCase := range testCases {
		w := testWrapper{
			Time: Time{
				Time:         testCase.Time,
				OutputLayout: testCase.Layout,
			},
		}
		b, err := xml.Marshal(w)
		require.NoError(t, err)

		assert.Equal(t, createXml(testCase.Text), string(b), testCase.Layout)
	}
}

func TestUniversalTime_UnmarshalBSONValue(t *testing.T) {
	for _, testCase := range testCases {
		w := testWrapper{}
		err := bson.UnmarshalExtJSON([]byte(createBson(testCase.Text)), true, &w)
		require.NoError(t, err)

		testTime(t, testCase.Text, testCase.Time, w.Time.Time)
	}
}

func TestUniversalTime_MarshalBSONValue(t *testing.T) {
	for _, testCase := range testCases {
		w := testWrapper{
			Time: Time{
				Time:         testCase.Time,
				OutputLayout: testCase.Layout,
			},
		}

		b, err := bson.MarshalExtJSON(w, true, true)
		require.NoError(t, err)

		assert.Equal(t, createBson(testCase.Text), string(b), testCase.Layout)
	}
}

func createJson(timeText string) string {
	return fmt.Sprintf(`{"time":"%s"}`, timeText)
}

func createXml(timeText string) string {
	return fmt.Sprintf(`<Wrapper><time>%s</time></Wrapper>`, timeText)
}

func createBson(timeText string) string {
	return fmt.Sprintf(`{"xmlname":{"space":"","local":""},"time":"%s","timenil":null}`, timeText)
}

func testTime(t *testing.T, text string, expected, actual time.Time) {
	assert.Equal(t, expected.Year(), actual.Year(), "year from \"%s\"", text)
	assert.Equal(t, expected.Month(), actual.Month(), "month from \"%s\"", text)
	assert.Equal(t, expected.Day(), actual.Day(), "day from \"%s\"", text)
	assert.Equal(t, expected.Hour(), actual.Hour(), "hour from \"%s\"", text)
	assert.Equal(t, expected.Minute(), actual.Minute(), "minute from \"%s\"", text)
	assert.Equal(t, expected.Second(), actual.Second(), "second from \"%s\"", text)
}

func TestUniversalTime_Split(t *testing.T) {
	o := ToUniversal(time.Date(2022, 7, 11, 10, 33, 56, 0, msk))
	d := o.Split()
	testTime(t, "", time.Date(2022, 7, 11, 0, 0, 0, 0, msk), d.Date.Time)
	testTime(t, "", time.Date(0, 1, 1, 10, 33, 56, 0, msk), d.Time.Time)

	b, err := json.Marshal(d)
	require.NoError(t, err)

	require.Equal(t, `{"Date":"2022-07-11","Time":"10:33:56"}`, string(b))
}
