package config

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type HttpServer struct {
	Address string `env:"HTTP_SERVER_ADDRESS" yaml:"address"`
	User    string `env:"HTTP_SERVER_USER" yaml:"user"`
	Pass    string `env:"HTTP_SERVER_PASS" yaml:"pass"`
}

type AppConfig struct {
	PathConfig  string       `env:"PATH_CONFIG"`
	ServiceName string       `env:"SERVICE_NAME" yaml:"serviceName"`
	Debug       bool         `env:"DEBUG" yaml:"debug"`
	RabbitMQ    rabbitConfig `yaml:"rabbitMQ"`
	Queues      queues       `yaml:"queues"`
	PGStorage   pgStorage    `yaml:"pgStorage"`
	HTTPServer  HttpServer   `yaml:"http_server"`
}

type queues struct {
	InfoMessage string `env:"INFO_MESSAGES_QUERY" yaml:"infoMessages"`
}

type rabbitConfig struct {
	Host     string `env:"RABBITMQ_HOST" yaml:"host"`
	Port     string `env:"RABBITMQ_PORT" yaml:"port"`
	Username string `env:"RABBITMQ_USERNAME" yaml:"username"`
	Password string `env:"RABBITMQ_PASSWORD" yaml:"password"`
	Path     string `env:"RABBITMQ_PATH" yaml:"path"`
}

// PGStorage configures pg repository.
type pgStorage struct {
	Host string `env:"STORAGE_HOST" yaml:"host"`
	Port int    `env:"STORAGE_PORT" yaml:"port"`
	User string `env:"STORAGE_PG_USER" yaml:"user"`
	Pass string `env:"STORAGE_PASS" yaml:"pass"`
	DB   string `env:"STORAGE_DB"   yaml:"db"`
}

// ReadYamlConfig reads the YAML configuration file and stores its contents in the AppConfig struct.
func (c *AppConfig) ReadYamlConfig(pathFile string) error {
	open, err := os.Open(pathFile)
	if err != nil {
		return err
	}
	defer func(open *os.File) {
		err = open.Close()
		if err != nil {
			log.Println(err)
		}
	}(open)
	if err = yaml.NewDecoder(open).Decode(&c); err != nil { //nolint:typecheck
		return err
	}
	return nil
}

// ReadEnvConfig reads the environment variables and stores their contents in the AppConfig struct.
func (c *AppConfig) ReadEnvConfig() error {
	if err := env.Parse(c); err != nil { //nolint:typecheck
		return err
	}
	return nil
}

// Validate validates the AppConfig struct.
func (c *AppConfig) Validate() error {
	if c.HTTPServer.Address == "" {
		return fmt.Errorf("address service is not set")
	}
	if c.PGStorage.Host == "" {
		return fmt.Errorf("pg repository host is not set")
	}
	// TODO: InfoMessages
	//if c.Queues.InfoMessages == "" {
	//	return fmt.Errorf("info messages query is not set")
	//}
	return nil
}
