package utils

import (
	"encoding/json"
	"os"
)

func ReadJson(f string, v interface{}) error {
	jsonFile, err := os.ReadFile(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonFile, v)
	if err != nil {
		return err
	}

	return nil
}

func WriteJson(f string, v interface{}) error {
	jsonFile, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = os.WriteFile(f, jsonFile, 0644)
	if err != nil {
		return err
	}

	return nil
}
