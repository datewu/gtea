package utils

import (
	"bytes"
	"encoding/gob"
)

// RedisMarshal use gob encode
func RedisMarshal(b any) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// RedisUnMarshal use gob decode
func RedisUnMarshal(data []byte, b any) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(b)
}
