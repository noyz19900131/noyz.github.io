package cfg

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port uint `yaml:"port"`
}
type DbConfig struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

type JwtConfig struct {
	Secret string `yaml:"secret"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	Db     DbConfig     `yaml:"db"`
	Jwt    JwtConfig    `yaml:"jwt"`
}

var CFG = loadConfig("cfg/config.yml")

func loadConfig(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil
	}

	return &config
}
