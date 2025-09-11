package gokafka

import (
	"encoding/json"
	"strings"
)

func (m Message) GetHeader(key string) (value []byte) {
	for _, header := range m.Headers {
		if strings.EqualFold(header.Key, key) {
			copyValue := make([]byte, len(header.Value))
			copy(copyValue, header.Value)

			return copyValue
		}
	}

	return nil
}

func NewJSONMessage(key string, data any) (Message, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return Message{}, err
	}

	return Message{
		Key:   []byte(key),
		Value: jsonBytes,
	}, nil
}
