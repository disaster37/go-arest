package client

import (
	"encoding/json"
	"errors"
)

// CheckMode permit to check that mode provided exist
func CheckMode(mode string) (err error) {
	if mode != OUTPUT && mode != INPUT && mode != INPUT_PULLUP {
		err = errors.New("Mode must be OUTPUT, INPUT or INPUT_PULLUP")
	}

	return err
}

// CheckLevel permit to check that level provided exist
func CheckLevel(level int) (err error) {
	if level != HIGH && level != LOW {
		err = errors.New("Level must be HIGH or LOW")
	}

	return err
}

// Unmarshal Parse JSON string to string map
func Unmarshal(msg []byte) (data map[string]interface{}, err error) {

	err = json.Unmarshal(msg, &data)

	return data, err
}
