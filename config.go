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

func (a *app) saveConfig() error {
	data, err := json.Marshal(a.config)
	if err != nil {
		return err
	}
	os.WriteFile(a.configPath, []byte(data), 0644)
	return nil
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
