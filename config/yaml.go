package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var filename = "config.yml"

type Config struct {
	Providers map[string]string `yaml:"providers"`
	Wallets   map[string]string `yaml:"wallets"`
}

func saveConfig(config *Config) error {
	file, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, file, 0644)
}

func LoadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// 파일이 없으면 빈 Config 반환
			return Config{Providers: make(map[string]string)}, nil
		}
		return config, err
	}
	err = yaml.Unmarshal(file, &config)
	return config, err
}

func AddProvider(config *Config, key, value string) error {
	if config.Providers == nil {
		config.Providers = make(map[string]string)
	}
	config.Providers[key] = value
	err := saveConfig(config)

	if err != nil {
		return err
	}
	return nil
}

func DeleteProvider(config *Config, key string) error {
	delete(config.Providers, key)
	err := saveConfig(config)

	if err != nil {
		return err
	}
	return nil
}
