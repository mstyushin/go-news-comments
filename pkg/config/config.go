package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	VersionDate = "0001-01-01 00:00:00"
	Version     = "dev"
	Hash        = "COMMIT ID"
)

// DefaultConfig возвращает конфигурацию по умолчанию.
// будет использоваться в случае когда config.yaml не найден
func DefaultConfig() *Config {
	return &Config{
		HttpPort:     8082,
		DBConnString: "postgres://postgres@localhost:5433/comments?sslmode=disable",
	}
}

// конфигурация приложения, подразумевается yaml-формат
type Config struct {
	HttpPort     int    `yaml:"http_port"`
	DBConnString string `yaml:"db_conn_string,omitempty"`
}

func VersionString() string {
	return fmt.Sprintf("Version: %s Commit: %s BuildDate: %s", Version, Hash, VersionDate)
}

func New() (*Config, error) {
	var config *Config

	var configPath string
	var printConfig bool
	var printVersion bool
	flag.StringVar(&configPath, "config", "./config.yaml", "path to a YAML config file")
	flag.BoolVar(&printConfig, "print-config", false, "print loaded config")
	flag.BoolVar(&printVersion, "version", false, "print build version")
	flag.Parse()

	if printVersion {
		fmt.Println(VersionString())
		return nil, nil
	}

	f, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("not found config file at %s, using defaults\n", configPath)
		config = DefaultConfig()
	} else {
		log.Printf("reading config at %s\n", configPath)
	}

	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return nil, err
	}

	if printConfig {
		yamlData, _ := yaml.Marshal(&config)
		fmt.Println()
		fmt.Println(string(yamlData))
		return nil, nil
	}

	return config, nil
}
