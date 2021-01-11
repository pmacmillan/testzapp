package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config ...
type Config struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func loadConfig(file string) (*Config, error) {
	var config Config

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
