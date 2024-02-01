package configs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	SourceKey string `yaml:"gemini-keys"`
}

func GetConfig(configPath string) Config {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("yamlFile.Get err #%v ", err))
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(fmt.Sprintf("Unmarshal: %v", err))
	}
	return config
}
