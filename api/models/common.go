package models

import (
	"encoding/json"
	"fmt"
)

/*
Serializes a struct into a string
*/
func Stringify(from interface{}) (string, error) {
	stringified, err := json.Marshal(from)
	if err != nil {
		return "", err
	}
	return string(stringified), nil
}

/*
Deserialize a string to a struct
*/
// Proxy
func Structify[T interface{}](from string, to T) error {
	return StructifyBytes([]byte(from), to)
}

/*
Deserialize a byte array into a struct
*/
func StructifyBytes[T interface{}](from []byte, to T) error {
	if len(from) == 0 {
		return fmt.Errorf("empty input string")
	}
	err := json.Unmarshal(from, to)
	return err
}
