package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var defaultServerSettings = Config{
	Logfile:    "messages.log",
	TelnetPort: "4444",
	HTTPPort:   "8080",
}

type Config struct {
	Logfile    string `yaml:"logfile"`
	TelnetPort string `yaml:"telnet_port"`
	HTTPPort   string `yaml:"http_port"`
}

func (c *Config) Validate() error {
	if c.Logfile == "" {
		return fmt.Errorf("logfile entry is required in configuration file")
	}

	if c.HTTPPort == "" {
		return fmt.Errorf("http_port is required in configuration file")
	}

	if c.TelnetPort == "" {
		return fmt.Errorf("telnet_port is required in configuration file")
	}
	return nil
}

func New(filename string) (*Config, error) {
	if filename == "" {
		log.Printf("[INFO] Default server config will be used")
		return &defaultServerSettings, nil
	}
	buff, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(buff, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid server config: %v", err)
	}
	return &cfg, nil
}
