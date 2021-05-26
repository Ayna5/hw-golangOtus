package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Config При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	DB     DBConf
	Server Server
}

type LoggerConf struct {
	Level string
	Path  string
}

type Server struct {
	HTTP string
	Grpc string
}

type DBConf struct {
	User     string
	Password string
	Host     string
	Port     uint64
	Name     string
}

func NewConfig(path string) (Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return Config{}, fmt.Errorf("cannot decode file: %w", err)
	}
	return config, nil
}
