package main

import (
	"encoding/json"
	"os"
)

//load config

// config struct
type config struct {
	ShowComplete bool `json:"show-complete"`
}

func (a *app) loadConfig() (*config, error) {
	data, err := os.ReadFile(a.configPath)
	if err != nil {
		return &config{}, err
	}
	var c config
	err = json.Unmarshal(data, &c)
	if err != nil {
		return &config{}, err
	}
	return &c, nil
}
