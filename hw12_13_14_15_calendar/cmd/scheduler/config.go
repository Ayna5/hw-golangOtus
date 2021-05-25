package main

import (
	"fmt"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/configs"
	"github.com/BurntSushi/toml"
)

func NewSheduler(path string) (configs.Sheduler, error) {
	var sheduler configs.Sheduler
	if _, err := toml.DecodeFile(path, &sheduler); err != nil {
		return configs.Sheduler{}, fmt.Errorf("cannot decode file: %w", err)
	}
	return sheduler, nil
}
