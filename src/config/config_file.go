package config

import (
	"encoding/json"
	"os"
)

func readConfigFile(path string) (Config, error) {
	var c Config

	content, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(content, &c)
	return c, err
}

func saveConfigFile(c Config, path string) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}
