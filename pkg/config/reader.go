package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ReadConfig() (*Config, error) {
	file, err := os.Open("config.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	err = decoder.Decode(&config)

	if err == nil {
		err = config.validate()
	}
	if err != nil {
		return nil, err
	}

	return &config, nil
}
