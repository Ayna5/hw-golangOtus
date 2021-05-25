package main

import (
	"fmt"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/configs"
	"github.com/BurntSushi/toml"
)

func NewConfig(path string) (configs.Config, error) {
	var config configs.Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return configs.Config{}, fmt.Errorf("cannot decode file: %w", err)
	}
	return config, nil
}
