package main

import (
	"fmt"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/configs"
	"github.com/BurntSushi/toml"
)

func NewSender(path string) (configs.Sender, error) {
	var sender configs.Sender
	if _, err := toml.DecodeFile(path, &sender); err != nil {
		return configs.Sender{}, fmt.Errorf("cannot decode file: %w", err)
	}
	return sender, nil
}
