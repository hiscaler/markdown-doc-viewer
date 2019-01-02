package config

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	DocumentDir string
	Port         int
}

func NewConfig() *Config {
	config := &Config{}

	if file, err := os.Open("./config/config.json"); err != nil {
		config.DocumentDir = ""
	} else {
		defer file.Close()
		jsonByte, err := ioutil.ReadAll(file)
		if err == nil {
			json.Unmarshal(jsonByte, &config)
		}
	}
	if len(config.DocumentDir) <= 0 {
		config.DocumentDir = "./docs"
	}

	return config
}
