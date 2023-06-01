package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	APIKey   string   `json:"apiKey"`
	CX       string   `json:"cx"`
	Keywords []string `json:"keywords"`
}

func ReadConfig() (Config, error) {
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
